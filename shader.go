package main

import "github.com/rajveermalviya/go-webgpu/wgpu"

type Shader struct {
	Module   *wgpu.ShaderModule
	Pipeline *wgpu.ComputePipeline
	Manager  *Manager
	Name     string
}

func (m *Manager) NewShader(name, entrypoint, code string) (*Shader, error) {
	var err error
	succsessful := false
	f := &Shader{Manager: m, Name: name}
	defer func() {
		if !succsessful {
			f.Release()
		}
	}()
	// Create shader module
	f.Module, err = m.Device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
		Label: name + "[shader]",
		WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
			Code: code,
		},
	})
	if err != nil {
		return nil, err
	}

	// Create pipeline
	f.Pipeline, err = m.Device.CreateComputePipeline(&wgpu.ComputePipelineDescriptor{
		Label: name + "[pipeline]",
		Compute: wgpu.ProgrammableStageDescriptor{
			Module:     f.Module,
			EntryPoint: entrypoint,
		},
	})
	if err != nil {
		return nil, err
	}

	succsessful = true
	return f, nil
}

func (f *Shader) Release() {
	if f.Pipeline != nil {
		f.Pipeline.Release()
	}
	if f.Module != nil {
		f.Module.Release()
	}
}
