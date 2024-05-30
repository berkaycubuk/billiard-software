import api from "./APIService";

const overview = async (id) => {
  const response = await api.get('/v1/payment/overview/' + id, true);
  return response;
}

const init = async (id, method) => {
  const response = await api.post('/v1/payment/init', true, {
    order_id: id,
    method: method,
  });
  return response;
}

const accept = async (id) => {
  const response = await api.post('/v1/payment/accept', true, {
    order_id: id,
    method: 1,
  });
  return response;
}

const status = async (id, method) => {
  const response = await api.post('/v1/payment/status', true, {
    order_id: id,
    method: method,
  });
  return response;
}

export default {
  overview,
  init,
  status,
  accept,
};
