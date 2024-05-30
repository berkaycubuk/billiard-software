import api from "./APIService";

const kiosk = async (start = null, end = null) => {
  let routeString = '/v1/reports/kiosk';
  let started = false;
  
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
  return response;
}

const tables = async (start = null, end = null) => {
  let routeString = '/v1/reports/table-sales/all';
  let started = false;
  
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
  return response;
}

const subscriptions = async (start = null, end = null) => {
  let routeString = '/v1/reports/subscriptions';
  let started = false;
  
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
  return response;
}

const total = async (start = null, end = null) => {
  let routeString = '/v1/reports/total';
  let started = false;
  
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
  return response;
}

const games = async (start = null, end = null, page = null) => {
  let routeString = '/v1/reports/games';
  let started = false;
  
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
  
  if (page != null) {
    if (started) {
      routeString += '&page=' + page;
    } else {
      routeString += '?page=' + page;
    }
    started = true;
  }

  const response = await api.get(routeString, true);
  return response;
}

const user = async (id, start = null, end = null) => {
  let routeString = '/v1/reports/users/' + id;
  let started = false;
  
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
  return response;
}

const users = async (start = null, end = null, user = null, page = null) => {
  let routeString = '/v1/reports/users';
  let started = false;
  
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
  
  if (user != null) {
    if (started) {
      routeString += '&user=' + user;
    } else {
      routeString += '?user=' + user;
    }
    started = true;
  }
  
  if (page != null) {
    if (started) {
      routeString += '&page=' + page;
    } else {
      routeString += '?page=' + page;
    }
    started = true;
  }

  const response = await api.get(routeString, true);
  return response;
}

export default {
  kiosk,
  tables,
  subscriptions,
  total,
  games,
  users,
  user,
};
