import api from "./APIService";

const all = async () => {
  const response = await api.get('/v1/table/all', true);
  return response;
}

const get = async (id) => {
  const response = await api.get('/v1/table/' + id, true);
  return response;
}

const join = async (id) => {
  const response = await api.post('/v1/table/join', true, {
    id
  });
  return response;
}

const del = async (id) => {
  const response = await api.post('/v1/table/delete', true, {
    id
  });
  return response;
}

const create = async (name) => {
  const response = await api.post('/v1/table/new', true, {
    name
  });
  return response;
}

const update = async (id, name) => {
  const response = await api.post('/v1/table/update', true, {
    id,
    name
  });
  return response;
}

const joinAsGuest = async (id, name) => {
  const response = await api.post('/v1/table/join-as-guest', true, {
    id,
    name
  });
  return response;
}

const transferGuest = async (id, newTableId, name) => {
  const response = await api.post('/v1/table/transfer-as-guest', true, {
    id,
    new_table_id: newTableId,
    name,
  });
  return response;
}

const transfer = async (id, newTableId) => {
  const response = await api.post('/v1/table/transfer', true, {
    id,
    new_table_id: newTableId,
  });
  return response;
}

const pause = async (id) => {
  const response = await api.post('/v1/table/pause', true, {
    id
  });
  return response;
}

const pauseAsGuest = async (id, name) => {
  const response = await api.post('/v1/table/pause-as-guest', true, {
    id,
    name
  });
  return response;
}

const unpause = async (id) => {
  const response = await api.post('/v1/table/unpause', true, {
    id
  });
  return response;
}

const unpauseAsGuest = async (id, name) => {
  const response = await api.post('/v1/table/unpause-as-guest', true, {
    id,
    name
  });
  return response;
}

const leaveAsGuest = async (id, name) => {
  const response = await api.post('/v1/table/leave-as-guest', true, {
    id,
    name
  });
  return response;
}

const leave = async (id) => {
  const response = await api.post('/v1/table/leave', true, {
    id
  });
  return response;
}

export default {
  all,
  get,
  join,
  del,
  leave,
  leaveAsGuest,
  pause,
  create,
  update,
  pauseAsGuest,
  unpause,
  unpauseAsGuest,
  joinAsGuest,
  transfer,
  transferGuest,
};
