import api from "./APIService";

const all = async () => {
  const response = await api.get('/v1/role/all', true);
  return response.data.roles;
}

const get = async (id) => {
  const response = await api.get('/v1/role/' + id, true);
  return response.data.role;
}

const update = async (id, name) => {
  const response = await api.post('/v1/role/update', true, {
    id,
    name,
  });
  return response.data.role;
}

const create = async (name) => {
  const response = await api.post('/v1/role/new', true, {
    name,
  });
  return response.data.success;
}

const del = async (id) => {
  const response = await api.post('/v1/role/delete', true, {
    id
  });
  return response.data.success;
}

export default {
  all,
  get,
  del,
  update,
  create,
};
