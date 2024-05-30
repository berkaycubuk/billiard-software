import api from "./APIService";

const buy = async (items) => {
  const response = await api.post('/v1/shop/buy', true, {
    items,
  });
  return response;
}

const addToGuestAccount = async (game_user_name, items) => {
  const response = await api.post('/v1/shop/add-to-guest-account', true, {
    game_user_name,
    items,
  });
  return response;
}

const addToAccount = async (items) => {
  const response = await api.post('/v1/shop/add-to-account', true, {
    items,
  });
  return response;
}

export default {
  buy,
  addToAccount,
  addToGuestAccount,
};
