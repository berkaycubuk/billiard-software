import API from "./APIService";

const all = async (user = null, page = null) => {
  let string = '/v1/user/all';
  let started = false;

  if (page != null) {
    string += '?page=' + page;

    started = true;
  }

  if (user != null) {
    if (started) {
      string += '&user=' + user;
    } else {
      string += '?user=' + user;
      started = true;
    }
  }

  const response = await API.get(string, true);
  
  return response.data.data;
}

const allWithoutPagination = async () => {
  let string = '/v1/user/all-without-pagination';
  const response = await API.get(string, true);
  
  return response.data.data;
}

const get = async (id) => {
  const response = await API.get('/v1/user/' + id, true);

  if (response.data.success === false) {
    return null;
  }
  
  return response.data.user;
}

const update = async (id, name, surname, email, phone, role) => {
  const response = await API.post('/v1/user/update', true, {
    id,
    name,
    surname,
    email,
    phone,
	  role,
  });
  return response.data.success;
}

const getProfile = async () => {
  try {
    const response = await API.get('/v1/user/profile', true)
    return response
  } catch (err) {
    throw err
  }
}

const myActiveSubscription = async () => {
  const response = await API.get('/v1/user/my-active-subscription', true);
  return response;
}

const mySubscriptions = async () => {
  const response = await API.get('/v1/user/my-subscriptions', true);
  return response;
}

const getGuests = async () => {
  const response = await API.get('/v1/user/guests', true);
  return response;
}

const del = async (id) => {
  const response = await API.post('/v1/user/delete/' + id, true, {});
  return response;
}

const create = async (name, surname, email, phone, emailVerified, role) => {
  try {
    const response = await API.post('/v1/user/create', true, {
      name,
      surname,
      email,
      phone,
      email_verified: emailVerified,
	  role,
    });

    return response.data;
  } catch(error) {
    console.log(error)
    throw error;
  }
}

export default {
  all,
  get,
  del,
  create,
  getProfile,
  getGuests,
  mySubscriptions,
  myActiveSubscription,
  update,
  allWithoutPagination,
};
