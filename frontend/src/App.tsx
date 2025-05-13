import { task } from "./entities/task.ts";
import { ToDoList } from "./components/ToDoList.tsx";
import { Header } from "./components/Header.tsx";

let task1 = new task("1", "Task 1: Complete project", "Medium", "Active");
let task2 = new task("2", "Task 2: Review code", "High", "Late", "Check all pull requests",
    new Date("2023-12-17T03:24:00"));

function App() {
    return (
        <div className="layout">
            <Header/>
            <div className="page-content">
                <ToDoList tasks={[task1, task2]} />
            </div>
        </div>
    )
}

export default App;