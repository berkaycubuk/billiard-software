// import APIService from "./APIService";

// const loggerService = new APIService();

const error = async (message) => {
  console.error(message);
}

const log = async (message) => {
  console.log(message);
}

export default {
  error,
  log
};
