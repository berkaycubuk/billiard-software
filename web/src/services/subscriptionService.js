import api from "./APIService";

const mySubscriptions = async () => {
  const response = await api.get('/v1/user/my-subscriptions', true);
  return response;
}

const all = async () => {
  const response = await api.get('/v1/subscription/all', true);

  if (response.data.success === false) {
    return null;
  }

  return response.data.subscriptions;
}

const buy = async (id) => {
  const response = await api.post('/v1/subscription/buy', true, {
    id: id,
  });

  return response;
}

const get = async (id) => {
  const response = await api.get('/v1/subscription/' + id, true);
  return response.data.subscription;
}

const update = async (id, name, price, hours, role, hidden) => {
  const response = await api.post('/v1/subscription/update', true, {
    id,
    name,
    price,
    hours,
	role,
	hidden,
  });
  return response.data.success;
}

const add = async (name, price, hours, role, hidden) => {
  const response = await api.post('/v1/subscription/new', true, {
    name,
    price,
    hours,
  	role,
	hidden,
  });
  return response.data.success;
}

const pause = async (id) => {
  const response = await api.post('/v1/subscription/pause', true, {
    id,
  });
  return response.data.success;
}

const unpause = async (id) => {
  const response = await api.post('/v1/subscription/unpause', true, {
    id,
  });
  return response.data.success;
}

const del = async (id) => {
  try {
    const response = await api.post('/v1/subscription/delete', true, {
      id,
    });
    return response.data.success;
  } catch (err) {
    throw err
  }
}

const delUserSub = async (id) => {
  try {
    const response = await api.post('/v1/subscription/delete-user-subscription', true, {
      id,
    });
    return response.data.success;
  } catch (err) {
    throw err
  }
}

const addUserSub = async (user_id, sub_id) => {
  try {
    const response = await api.post('/v1/subscription/add-user-subscription', true, {
      user_id,
      sub_id,
    });
    return response.data.success;
  } catch (err) {
    throw err
  }
}

const pauseUserSub = async (user_id, sub_id) => {
  try {
    const response = await api.post('/v1/subscription/pause-user-subscription', true, {
      user_id,
      sub_id,
    });
    return response.data.success;
  } catch (err) {
    throw err
  }
}

const hoursToHumanReadable = (hours) => {
  if (hours <= 24) {
    return `${hours} hours`;
  }

  return `${hours / 24} days`;
}

export default {
  all,
  buy,
  get,
  add,
  del,
  update,
  pause,
  unpause,
  mySubscriptions,
  hoursToHumanReadable,
  delUserSub,
  addUserSub,
  pauseUserSub,
};
