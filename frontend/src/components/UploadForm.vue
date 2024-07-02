<template>
  <div class="container mx-auto p-4 font-sans">
    <h2 class="text-3xl font-bold mb-6">Upload VRP Files</h2>
    <input type="file" @change="handleFileUpload" multiple />
    <button
      @click="submitFiles"
      class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
    >
      Upload
    </button>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      files: [],
    };
  },
  methods: {
    handleFileUpload(event) {
      this.files = event.target.files;
    },
    async submitFiles() {
      const formData = new FormData();
      for (let file of this.files) {
        formData.append("files", file);
      }

      try {
        console.log("Test");
        const response = await axios.post(
          "http://localhost:8080/upload",
          formData
        );

        console.log(response);
        console.log("Files uploaded");
        alert("Files uploaded successfully");
      } catch (error) {
        console.log(error);
        alert("Error uploading files");
      }
    },
  },
};
</script>
