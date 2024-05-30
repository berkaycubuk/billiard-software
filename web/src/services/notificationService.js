import api from "./APIService";

const active = async () => {
  let string = '/v1/notification/active';

  const response = await api.get(string, true);
  
  return response.data.notifications;
}

const game = async () => {
  let string = '/v1/notification/game';

  const response = await api.get(string, true);
  
  return response.data.notifications;
}

const mark = async (id) => {
  let string = '/v1/notification/mark';

  const response = await api.post(string, true, {
    id,
  });
  
  return response.data.data;
}

const genNot = async () => {
  let string = '/v1/notification/gen';

  const response = await api.post(string, true, {});
  
  return response.data.data;
}

const create = async (message) => {
  let string = '/v1/notification/create';

  const response = await api.post(string, true, {
    message,
  });
  
  return response.data.data;
}

export default {
  active,
  mark,
  genNot,
  game,
  create,
};
