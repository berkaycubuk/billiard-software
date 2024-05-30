import api from "./APIService";

const all = async (page = null, start = null, end = null, method = null, zeroOrders = false) => {
  let routeString = '/v1/order/all';
  let started = false;

  if (page != null) {
    routeString += '?page=' + page;

    started = true;
  }
  
  if (start != null) {
    if (started) {
      routeString += '&start=' + start;
    } else {
      routeString += '?start=' + start;
    }
    started = true;
  }
  
  if (end != null) {
    if (started) {
      routeString += '&end=' + end;
    } else {
      routeString += '?end=' + end;
    }
    started = true;
  }

  if (method != null) {
    if (started) {
      routeString += '&method=' + method;
    } else {
      routeString += '?method=' + method;
    }
    started = true;
  }

  if (zeroOrders == true || zeroOrders == "true") {
    if (started) {
      routeString += '&zeroOrders=true';
    } else {
      routeString += '?zeroOrders=true';
    }
    started = true;
  }

  const response = await api.get(routeString, true);
  return response.data.data;
}

const myOldOrders = async (page = null, start = null, end = null) => {
  let routeString = '/v1/order/my-old-orders';
  let started = false;

  if (page != null) {
    routeString += '?page=' + page;

    started = true;
  }
  
  if (start != null) {
    if (started) {
      routeString += '&start=' + start;
    } else {
      routeString += '?start=' + start;
    }
    started = true;
  }
  
  if (end != null) {
    if (started) {
      routeString += '&end=' + end;
    } else {
      routeString += '?end=' + end;
    }
    started = true;
  }

  const response = await api.get(routeString, true);
  return response.data.data;
}

const myWaitingOrders = async () => {
  const response = await api.get('/v1/order/my-waiting-orders', true);
  return response.data.data;
}

const userOrders = async (id) => {
  const response = await api.get('/v1/order/user-orders/' + id, true);
  return response.data.data;
}

const capture = async (ref) => {
  try {
    const response = await api.get('/v1/payment/capture/' + ref, true);
    return response.data;
  } catch(err) {
    throw err;
  }
}

const approve = async (id) => {
  const response = await api.post('/v1/order/approve/' + id, true);
  return response.data.success;
}

const applyDiscount = async (id, amount) => {
  const response = await api.post('/v1/order/apply-discount', true, {
    id,
    amount,
  });
  return response.data.success;
}

const get = async (id) => {
  const response = await api.get('/v1/order/' + id, true);
  return response.data.order;
}

// delete is forbidden so come up with this for now
const del = async (id) => {
  const response = await api.post('/v1/order/delete/' + id, true);
  return response.data.success;
}

const cancel = async (id) => {
  const response = await api.post('/v1/order/cancel/' + id, true);
  return response.data.success;
}

const actionString = (code) => {
  if (code === 1) {
    return "actions.init";
  } else if (code === 2) {
    return "actions.cancel";
  } else if (code === 3) {
    return "actions.pay";
  } else if (code === 4) {
    return "actions.transfer";
  } else if (code === 5) {
    return "actions.approve";
  }
}

const methodString = (code) => {
  if (code === 1) {
    return "methods.physical";
  } else if (code === 2) {
    return "methods.vipps";
  } else if (code === 3) {
    return "methods.system";
  }
}

const statusString = (code) => {
  if (code === 1) {
    return "statuses.waiting";
  } else if (code === 2) {
    return "statuses.canceled";
  } else if (code === 3) {
    return "statuses.paid";
  } else if (code === 4) {
    return "statuses.deleted";
  } else if (code === 5) {
    return "statuses.pay_later";
  } else if (code === 6) {
    return "statuses.transferred";
  } else if (code === 7) {
    return "statuses.approved";
  } else if (code === 8) {
    return "statuses.paid";
  }

  return "null";
}

export default {
  all,
  get,
  del,
  cancel,
  approve,
  userOrders,
  statusString,
  actionString,
  applyDiscount,
  capture,
  methodString,
  myWaitingOrders,
  myOldOrders,
};
