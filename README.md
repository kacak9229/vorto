## How to Run Frontend and Backend

### Using Docker

1. Run the following command:
   docker-compose up --build

   1. Access localhost:3000 -> Frontend
   2. Access localhost:8000 -> Backend

### Manually

#### Backend

1. Navigate to the backend directory:
   cd backend

2. Run the backend:
   go run main.go

3. The backend should now be running on:
   http://localhost:8080

#### Frontend

1. Navigate to the frontend directory:
   cd frontend

2. Install the dependencies:
   npm install

3. Run the frontend:
   npm run serve

4. The frontend should now be running on:
   http://localhost:8081

## Steps to Use

1. **Upload a File**:

- Upload one or multiple files.

2. **Click Upload**:

- Click the "Upload" button to send the files.

3. **Fetch Problems**:

- Fetch the problems to start solving.

## Algorithm Explanation

The `solveVRP` function in the backend uses a combination of a greedy algorithm and the nearest neighbor approach.

### Initial Approach

- Initially, I used a sequential approach to assign the loads to the driver, but it had higher costs.

### Optimization

- To optimize, I sort the loads by Euclidean distance to form clusters, reducing costs.
- The nearest neighbor algorithm starts at a point and then looks for the closest next point.
- This heuristic method makes the best immediate choice at each step (greedy strategy).

### Further Optimizations

- There are many other algorithms that can further optimize the solution. Each has its pros and cons.
- Examples include:
- Clark-Wright Algorithm
- Metaheuristic approaches like Genetic Algorithms (GA), Tabu Search, and Ant Colony Optimization.

### Simplified Algorithm Explanation

- The algorithm assigns loads to drivers by always picking the next closest load to save time.
- If adding a load makes the driver work more than 12 hours, it finishes with that driver and starts with a new one.
- This step-by-step method builds the solution by making the best immediate choice each
