<!-- FileUpload.vue -->
<template>
    <div>
      <file-pond
        ref="pond"
        name="file"
        accepted-file-types="image/jpeg, image/png"
        max-files="1"
        :server="{
            url: APIAddress,
            process: {
                url: '/v1/upload',
                headers: {
                    'Authorization': 'bearer ' + cookies.get('token'),
                }
            }
        }"
      />
    </div>
  </template>
  
<script>
    import vueFilePond from "vue-filepond";
  import 'filepond/dist/filepond.min.css';
    import { useCookies } from 'vue3-cookies';

  const FilePond = vueFilePond();
  
  export default {
    components: {
      FilePond,
    },
    setup() {
        const { cookies } = useCookies();
        return { cookies };
    },
    data() {
        return {
            APIAddress: import.meta.env.VITE_API_URL,
        }
    },
  };
</script>