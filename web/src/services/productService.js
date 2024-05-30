import api from "./APIService";

const all = async () => {
  const response = await api.get('/v1/product/all', true);
  return response;
}

const get = async (id) => {
  const response = await api.get('/v1/product/' + id, true);
  return response.data.product;
}

const del = async (id) => {
  const response = await api.post('/v1/product/delete', true, {
    id,
  });
  return response.data.success;
}

const update = async (id, name, order, price, image = null) => {
  const response = await api.post('/v1/product/update', true, {
    id,
    name,
    order,
    price,
    image: image === null ? 0 : image,
  });
  return response.data.success;
}

const updateOrder = async (id, direction) => {
  const response = await api.post('/v1/product/update-order', true, {
    id,
    direction,
  });
  return response.data.success;
}

const create = async (name, price, image = null) => {
  const response = await api.post('/v1/product/new', true, {
    name,
    price,
    image: image === null ? 0 : image,
  });
  return response;
}

export default {
  all,
  create,
  get,
  del,
  update,
  updateOrder,
};
