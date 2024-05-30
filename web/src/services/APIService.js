import router from '@/router';
import axios from 'axios';
import { useCookies } from 'vue3-cookies';
import { toast } from 'vue3-toastify';
import { $t } from '../config'

const baseURL = import.meta.env.VITE_API_URL || "http://127.0.0.1:4000";

function checkError(error) {
  if (error.code === "ERR_NETWORK") {
    router.push('/network-error');
  }

  if (error.response.data.message) {
    toast.error($t(error.response.data.message));
  } else {
    toast.error(error.code);
  }

  if (error.response.data.message == 'error.auth_required') {
    setTimeout(() => {
      router.push('/login')
    }, 1000)
  }

  console.error(error);
}

const get = async (endpoint, auth = false, config = {}) => {
  const { cookies } = useCookies();
  try {
    let newConfig;
    if (auth) {
      newConfig = {...config, headers: {
        'Authorization': 'bearer ' + cookies.get('token'),
      }}
    } else {
      newConfig = config
    }

    const response = await axios.get(baseURL + endpoint, newConfig);
    return response;
  } catch(error) {
    checkError(error);
    throw error
  }
}

const post = async (endpoint, auth = false, data = {}, config = {}) => {
  const { cookies } = useCookies();
  try {
    let newConfig;
    if (auth) {
      newConfig = {...config, headers: {
        'Authorization': 'bearer ' + cookies.get('token'),
      }}
    } else {
      newConfig = config
    }

    const response = await axios.post(baseURL + endpoint, data, newConfig);

    // fix this logic to pass error messages to services
    if (response.data.success === false) {
      console.log(response.data);
    }

    return response;
  } catch(error) {
    checkError(error);
    throw error
  }
}

export default {
  get,
  post,
}
