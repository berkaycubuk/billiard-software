export const APIURL = () => {
    return import.meta.env.VITE_API_URL;
}

export const UploadsURL = () => {
    return APIURL() + "/static/uploads";
}