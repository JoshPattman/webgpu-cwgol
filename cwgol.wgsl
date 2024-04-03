@group(0) @binding(0) var<storage, read_write> prev: array<f32>;
@group(0) @binding(1) var<storage, read_write> next: array<f32>;
@group(0) @binding(2) var<storage, read_write> shape: array<u32>;

@compute
@workgroup_size(16, 16, 1)
fn cwgol( @builtin(global_invocation_id) id: vec3<u32> ) {
    let ids = vec3<i32>(i32(id.x), i32(id.y), i32(id.z));
    if (!index_within_bounds(ids)) {
        return;
    }
    var count = 0.0;
    for (var dx: i32 = -1; dx <= 1; dx+=1) {
        for (var dy: i32 = -1; dy <= 1; dy+=1) {
            let did = vec3<i32>(ids.x + dx, ids.y + dy, ids.z);
            if ((dx == 0 && dy == 0) || !index_within_bounds(did)) {
                continue;
            }
            count += prev[get_index(did.x, did.y)];
        }
    }
    let i = get_index(ids.x, ids.y);
    let alive = prev[i] > 0.5;
    if (alive && count < 2.0) {
        next[i] = 0.0;
    } else if (alive && count > 3.0) {
        next[i] = 0.0;
    } else if (!alive && count == 3.0) {
        next[i] = 1.0;
    } else {
        next[i] = prev[i];
    }
}

// index within bounds function
fn index_within_bounds( id: vec3<i32> ) -> bool {
    return id.x >= 0 && id.y >= 0 && id.z >= 0 && id.x < i32(shape[0]) && id.y < i32(shape[1]) && id.z < i32(shape[2]);
}

// get the index in the array of an xy point
fn get_index( x: i32, y: i32 ) -> u32 {
    return u32(y * i32(shape[0]) + x);
}