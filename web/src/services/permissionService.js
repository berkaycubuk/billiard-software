import api from "./APIService";

const permission = async (name) => {
  const response = await api.get('/v1/user/permissions', true);

  if (response.data.success === false) {
    return null;
  }

  const permissions = response.data.permissions;
  if (permissions === null) {
    return null;
  }
  const found = permissions.find((permission) => permission.name === name);

  if (found === null) {
    return null;
  }

  return found.value;
}

export default {
  permission,
};
