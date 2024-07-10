<template>
  <div class="container mx-auto p-4 font-sans">
    <h2 class="text-3xl font-bold mb-6">VRP Results</h2>
    <button
      @click="fetchProblems"
      class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
    >
      Fetch Problems
    </button>

    <div v-if="problems.length" class="mt-8">
      <h3 class="text-2xl font-semibold mb-4">Problems</h3>
      <ul class="list-disc pl-5">
        <li
          v-for="problem in problems"
          :key="problem.id"
          class="cursor-pointer text-blue-500"
          @click="fetchResults(problem.id)"
        >
          Problem ID: {{ problem.id }} (Created at:
          {{ new Date(problem.created_at).toLocaleString() }})
        </li>
      </ul>
    </div>

    <div v-if="drivers.length" class="mt-8">
      <h3 class="text-2xl font-semibold mb-4">Drivers</h3>
      <div class="mb-4 flex justify-center pagination-container">
        <button
          @click="prevPage"
          :disabled="currentPage === 0"
          class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded-l"
        >
          Prev
        </button>
        <button
          v-for="page in pageRange"
          :key="page"
          @click="setPage(page - 1)"
          :class="[
            'px-4 py-2 mx-1 font-bold',
            page === currentPage + 1
              ? 'bg-blue-500 text-white'
              : 'bg-gray-300 hover:bg-gray-400 text-gray-800',
          ]"
        >
          {{ page }}
        </button>
        <button
          @click="nextPage"
          :disabled="currentPage >= pageCount - 1"
          class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded-r"
        >
          Next
        </button>
      </div>
      <table class="min-w-full bg-white border border-gray-200 rounded-lg">
        <thead>
          <tr
            class="w-full bg-gray-200 text-left text-gray-700 uppercase text-sm leading-normal"
          >
            <th class="py-3 px-4">Driver</th>
            <th class="py-3 px-4">Total Time (minutes)</th>
            <th class="py-3 px-4">Loads</th>
          </tr>
        </thead>
        <tbody class="text-gray-700 text-sm font-light">
          <tr
            v-for="(driver, index) in paginatedDrivers"
            :key="currentPage * pageSize + index"
            class="border-b border-gray-200 hover:bg-gray-100"
          >
            <td class="py-3 px-4">
              Driver {{ currentPage * pageSize + index + 1 }}
            </td>
            <td class="py-3 px-4">{{ Math.ceil(driver.TotalTime) }} minutes</td>
            <td class="py-3 px-4">
              <ul>
                <li v-for="load in driver.Loads" :key="load.LoadNumber">
                  Load {{ load.LoadNumber }}: ({{ load.PickupX }},
                  {{ load.PickupY }}) to ({{ load.DropoffX }},
                  {{ load.DropoffY }})
                </li>
              </ul>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="kpis" class="mt-8">
      <h3 class="text-2xl font-semibold mb-4">KPI Summary</h3>
      <table class="min-w-full bg-white border border-gray-200 rounded-lg">
        <thead>
          <tr
            class="w-full bg-gray-200 text-left text-gray-700 uppercase text-sm leading-normal"
          >
            <th class="py-3 px-4">Metric</th>
            <th class="py-3 px-4">Value</th>
          </tr>
        </thead>
        <tbody class="text-gray-700 text-sm font-light">
          <tr class="border-b border-gray-200 hover:bg-gray-100">
            <td class="py-3 px-4">Number of drivers</td>
            <td class="py-3 px-4">{{ kpis.number_of_drivers }}</td>
          </tr>
          <tr class="border-b border-gray-200 hover:bg-gray-100">
            <td class="py-3 px-4">Total cost</td>
            <td class="py-3 px-4">{{ formatCost(kpis.total_cost) }}</td>
          </tr>
          <tr class="border-b border-gray-200 hover:bg-gray-100">
            <td class="py-3 px-4">Mean cost</td>
            <td class="py-3 px-4">{{ formatCost(kpis.mean_cost) }}</td>
          </tr>
          <tr class="border-b border-gray-200 hover:bg-gray-100">
            <td class="py-3 px-4">Time taken to calculate</td>
            <td class="py-3 px-4">{{ kpis.time_taken }} seconds</td>
          </tr>
          <template
            v-for="fileName in Object.keys(kpis.drivers_per_file)"
            :key="fileName"
          >
            <tr class="bg-gray-100 font-semibold">
              <td colspan="2" class="py-3 px-4">{{ fileName }}</td>
            </tr>
            <tr class="border-b border-gray-200 hover:bg-gray-100">
              <td class="py-3 px-4">{{ fileName }} requires</td>
              <td class="py-3 px-4">
                {{ kpis.drivers_per_file[fileName] }} drivers
              </td>
            </tr>
            <tr class="border-b border-gray-200 hover:bg-gray-100">
              <td class="py-3 px-4">{{ fileName }} duration</td>
              <td class="py-3 px-4">
                {{ kpis.file_durations[fileName] }} seconds
              </td>
            </tr>
            <tr class="border-b border-gray-200 hover:bg-gray-100">
              <td class="py-3 px-4">{{ fileName }} cost</td>
              <td class="py-3 px-4">
                {{ formatCost(kpis.file_costs[fileName]) }}
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      problems: [],
      drivers: [],
      kpis: null,
      currentPage: 0,
      pageSize: 10,
    };
  },
  computed: {
    pageCount() {
      return Math.ceil(this.drivers.length / this.pageSize);
    },
    paginatedDrivers() {
      const start = this.currentPage * this.pageSize;
      const end = start + this.pageSize;
      return this.drivers.slice(start, end);
    },
    pageRange() {
      const rangeSize = 5;
      let start = Math.max(1, this.currentPage + 1 - Math.floor(rangeSize / 2));
      let end = Math.min(this.pageCount, start + rangeSize - 1);

      // Adjust start if end of pages range is beyond pageCount
      start = Math.max(1, end - rangeSize + 1);

      return Array.from({ length: end - start + 1 }, (v, i) => i + start);
    },
  },
  methods: {
    async fetchProblems() {
      try {
        const response = await axios.get("http://localhost:8080/problems");
        this.problems = response.data.problems;
      } catch (error) {
        alert("Error fetching problems");
      }
    },
    async fetchResults(problemId) {
      try {
        const response = await axios.get(
          `http://localhost:8080/results?problem_id=${problemId}`
        );
        this.drivers = response.data.drivers || [];
        this.kpis = response.data.kpis || null;

        console.log(this.kpis);
      } catch (error) {
        alert("Error fetching results");
      }
    },
    formatCost(value) {
      return "$ " + Math.ceil(value).toLocaleString();
    },
    nextPage() {
      if (this.currentPage < this.pageCount - 1) {
        this.currentPage++;
      }
    },
    prevPage() {
      if (this.currentPage > 0) {
        this.currentPage--;
      }
    },
    setPage(page) {
      this.currentPage = page;
    },
  },
};
</script>

<style scoped>
.pagination-container {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
