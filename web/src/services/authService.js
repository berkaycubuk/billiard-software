import api from "./APIService";
import { useCookies } from "vue3-cookies";

const { cookies } = useCookies();

const login = async (email, password, remember = false) => {
  try {
    const response = await api.post('/v1/auth/login', false, {
      email,
      password
    });

	if (remember == true) {
      cookies.set('token', response.data.token, '120d');
	} else {
      cookies.set('token', response.data.token, '30d');
	}
  } catch(e) {
    throw e
  }
}

const register = async (
  name,
  surname,
  email,
  phone,
  password,
  confirmPassword
) => {
  const response = await api.post('/v1/auth/register', false, {
    name,
    surname,
    email,
    phone,
    password,
    confirm_password: confirmPassword
  });

  return response;
}

const updatePassword = async (currentPassword, newPassword) => {
  const response = await api.post('/v1/auth/update-password', true, {
    current_password: currentPassword,
    new_password: newPassword,
  });

  return response;
}

const forgotPassword = async (email) => {
  const response = await api.post('/v1/auth/forgot-password', false, {
    email,
  });

  return response;
}

const passwordResetComplete = async (token, password) => {
  const response = await api.post('/v1/auth/password-reset-complete', false, {
    token,
    password,
  });

  return response;
}

const verify = async (token) => {
  const response = await api.post('/v1/auth/verify', false, {
    token,
  });

  return response;
}

export default {
  login,
  register,
  updatePassword,
  forgotPassword,
  passwordResetComplete,
  verify,
};
