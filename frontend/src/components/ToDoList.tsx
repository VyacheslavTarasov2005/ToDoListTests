import type { task } from "../entities/task.ts";
import { Task } from "./Task.tsx";

type Props = {
    tasks: task[];
}

export const ToDoList = ({ tasks }: Props) => {
    return(
        <div className='todolist'>
            <form className="get-form">
                <h1>Get All Tasks</h1>
                <select defaultValue="">
                    <option value="" disabled hidden>Choose sorting</option>
                    <option value="CreateAsc">Create Time Ascending</option>
                    <option value="CreateDesc">Create Time Descending</option>
                    <option value="PriorityAsc">Priority Ascending</option>
                    <option value="PriorityDesc">Priority Descending</option>
                    <option value="DeadlineAsc">Deadline Ascending</option>
                    <option value="DeadlineDesc">Deadline Descending</option>
                </select>
                <button>
                    Get
                </button>
            </form>
            <form className="create-task-form">
                <h1>Create Task</h1>
                <div>
                    <p>Name:</p>
                    <input/>
                </div>
                <div>
                    <p>Description:</p>
                    <input/>
                </div>
                <div>
                    <p>Deadline:</p>
                    <input type="date"/>
                </div>
                <div>
                    <p>Priority</p>
                    <select defaultValue="">
                        <option value="" disabled hidden>Choose priority</option>
                        <option value="Low">Low</option>
                        <option value="Medium">Medium</option>
                        <option value="High">High</option>
                        <option value="Critical">Critical</option>
                    </select>
                </div>
                <button>
                Create
                </button>
            </form>
            {tasks.map((task) => (
                <Task
                    key={task.id}
                    title={task.name}
                    description={task.description}
                    deadline={task.deadline}
                    priority={task.priority}
                    status={task.status}
                    isDone={task.isDone}
                />
            ))}
        </div>
    )
}