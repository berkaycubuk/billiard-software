import api from "./APIService";

const all = async () => {
  const response = await api.get('/v1/pricing/all', true);
  return response.data.pricings;
}

const get = async (id) => {
  const response = await api.get('/v1/pricing/' + id, true);
  return response.data.pricings;
}

const del = async (id) => {
  const response = await api.post('/v1/pricing/delete', true, {
    id,
  });
  return response.data.success;
}

const update = async (id, roleID, subscriptionID, playerCount, perMinute) => {
  const response = await api.post('/v1/pricing/update', true, {
    id,
    role_id: roleID,
    subscription_id: subscriptionID,
    player_count: playerCount,
    per_minute: perMinute,
  });
  return response.data.success;
}

const create = async (roleID, subscriptionID, playerCount, perMinute) => {
  const response = await api.post('/v1/pricing/create', true, {
    role_id: roleID,
    subscription_id: subscriptionID,
    player_count: playerCount,
    per_minute: perMinute,
  });
  return response;
}

export default {
  all,
  create,
  get,
  del,
  update,
};
