<template>
  <div id="app">
    <UploadForm />
    <GetResults />
  </div>
</template>

<script>
import UploadForm from "./components/UploadForm.vue";
import GetResults from "./components/GetResults.vue";

export default {
  components: {
    UploadForm,
    GetResults,
  },
};
</script>
