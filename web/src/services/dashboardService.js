import api from "./APIService";

const activeOrders = async () => {
  let routeString = '/v1/order/active-orders';
  const response = await api.get(routeString, true);
  return response;
}

const activeGames = async () => {
  let routeString = '/v1/dashboard/active-games';
  const response = await api.get(routeString, true);
  return response;
}

const kickUser = async (id) => {
  let routeString = '/v1/dashboard/kick-user';
  const response = await api.post(routeString, true, {id});
  return response;
}

const pauseUnpause = async (id) => {
  let routeString = '/v1/dashboard/pause-unpause';
  const response = await api.post(routeString, true, {id});
  return response;
}

const leaveTable = async (id) => {
  let routeString = '/v1/dashboard/leave-table';
  const response = await api.post(routeString, true, {id});
  return response;
}

export default {
  activeOrders,
  activeGames,
  pauseUnpause,
  leaveTable,
  kickUser,
};
