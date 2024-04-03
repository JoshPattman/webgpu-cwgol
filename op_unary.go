package main

import (
	"github.com/rajveermalviya/go-webgpu/wgpu"
)

type UnaryOp struct {
	WorkXBuffer  *wgpu.Buffer
	WorkYBuffer  *wgpu.Buffer
	ResultBuffer *wgpu.Buffer
	ShapeBuffer  *wgpu.Buffer
	BindGroup    *wgpu.BindGroup
	Function     *Shader
}

func (f *Shader) NewUnaryOp(arrLen uint64) (*UnaryOp, error) {
	var err error
	succsessful := false
	u := &UnaryOp{Function: f}
	defer func() {
		if !succsessful {
			u.Release()
		}
	}()

	// Create the buffers
	u.WorkXBuffer, err = f.Manager.Device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: f.Name + "[x buffer]",
		Size:  4 * arrLen,
		Usage: wgpu.BufferUsage_Storage | wgpu.BufferUsage_CopySrc | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return nil, err
	}

	u.WorkYBuffer, err = f.Manager.Device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: f.Name + "[y buffer]",
		Size:  4 * arrLen,
		Usage: wgpu.BufferUsage_Storage | wgpu.BufferUsage_CopySrc | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return nil, err
	}

	u.ResultBuffer, err = f.Manager.Device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: f.Name + "[mappable buffer]",
		Size:  4 * arrLen,
		Usage: wgpu.BufferUsage_MapRead | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return nil, err
	}

	u.ShapeBuffer, err = f.Manager.Device.CreateBuffer(&wgpu.BufferDescriptor{
		Label: f.Name + "[shape buffer]",
		Size:  4 * 3,
		Usage: wgpu.BufferUsage_Storage | wgpu.BufferUsage_CopySrc | wgpu.BufferUsage_CopyDst,
	})
	if err != nil {
		return nil, err
	}

	// Create the bind group to tell shader which buffers to use for computation
	u.BindGroup, err = f.Manager.Device.CreateBindGroup(&wgpu.BindGroupDescriptor{
		Label:  f.Name + "[bind group]",
		Layout: f.Pipeline.GetBindGroupLayout(0), // group(0)
		Entries: []wgpu.BindGroupEntry{
			{
				Binding: 0, // binding(0)
				Buffer:  u.WorkXBuffer,
				Size:    wgpu.WholeSize, // Not sure why, but this is needed
			},
			{
				Binding: 1, // binding(1)
				Buffer:  u.WorkYBuffer,
				Size:    wgpu.WholeSize, // Not sure why, but this is needed
			},
			{
				Binding: 2, // binding(2)
				Buffer:  u.ShapeBuffer,
				Size:    wgpu.WholeSize, // Not sure why, but this is needed
			},
		},
	})
	if err != nil {
		return nil, err
	}

	succsessful = true
	return u, nil
}

func (u *UnaryOp) Release() {
	if u.BindGroup != nil {
		u.BindGroup.Release()
	}
	if u.ResultBuffer != nil {
		u.ResultBuffer.Release()
	}
	if u.WorkXBuffer != nil {
		u.WorkXBuffer.Release()
	}
	if u.WorkYBuffer != nil {
		u.WorkYBuffer.Release()
	}
	if u.ShapeBuffer != nil {
		u.ShapeBuffer.Release()
	}
}

func (u *UnaryOp) Do(data []float32, workgroups, shape [3]uint32) ([]float32, error) {
	// Write the data
	u.Function.Manager.Device.GetQueue().WriteBuffer(u.WorkXBuffer, 0, wgpu.ToBytes(data))
	u.Function.Manager.Device.GetQueue().WriteBuffer(u.ShapeBuffer, 0, wgpu.ToBytes(shape[:]))

	// Create command encoder and do a pass
	encoder, err := u.Function.Manager.Device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
		Label: u.Function.Name + "[op encoder]",
	})
	if err != nil {
		return nil, err
	}
	defer encoder.Release()

	pass := encoder.BeginComputePass(&wgpu.ComputePassDescriptor{
		Label: u.Function.Name + "[op compute pass]",
	})

	pass.SetPipeline(u.Function.Pipeline)
	pass.SetBindGroup(0, u.BindGroup, nil)
	pass.DispatchWorkgroups(workgroups[0], workgroups[1], workgroups[2])
	err = pass.End()
	if err != nil {
		return nil, err
	}
	defer pass.Release()

	encoder.CopyBufferToBuffer(u.WorkYBuffer, 0, u.ResultBuffer, 0, u.ResultBuffer.GetSize())
	commandBuf, err := encoder.Finish(&wgpu.CommandBufferDescriptor{
		Label: u.Function.Name + "[op command buffer]",
	})
	if err != nil {
		return nil, err
	}
	defer commandBuf.Release()
	u.Function.Manager.Device.GetQueue().Submit(commandBuf)

	u.ResultBuffer.MapAsync(wgpu.MapMode_Read, 0, u.ResultBuffer.GetSize(), func(s wgpu.BufferMapAsyncStatus) {})
	u.Function.Manager.Device.Poll(true, nil)
	res := u.ResultBuffer.GetMappedRange(0, uint(u.ResultBuffer.GetSize()))
	resArr := wgpu.FromBytes[float32](res)
	result := make([]float32, len(resArr))
	copy(result, resArr)
	u.ResultBuffer.Unmap()
	return result, nil
}
