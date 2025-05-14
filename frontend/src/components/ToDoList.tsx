import { useEffect, useState } from "react";
import { Task } from "./Task";
import { task } from "../entities/task";
import { fetchTasks, createTask, deleteTask, toggleTaskStatus, updateTask } from "../api/api.ts";

export const ToDoList = () => {
    const [tasks, setTasks] = useState<task[]>([]);
    const [sorting, setSorting] = useState<string>("");
    const [editingTask, setEditingTask] = useState<task | null>(null);
    const [form, setForm] = useState({
        name: "",
        description: "",
        deadline: "",
        priority: "",
    });
    const [editForm, setEditForm] = useState({
        name: "",
        description: "",
        deadline: "",
        priority: "",
    });

    const loadTasks = async () => {
        try {
            const data = await fetchTasks(sorting as any);
            console.log("data:", data);
            const parsed = data.map((t: any) =>
                new task(
                    t.id,
                    t.name,
                    t.priority,
                    t.status,
                    t.createdAt,
                    t.changedAt,
                    t.description,
                    t.deadline ? new Date(t.deadline) : undefined
                )
            );
            setTasks(parsed);
            console.log("parsed:", parsed);
        } catch (err) {
            alert("Ошибка загрузки задач");
        }
    };

    useEffect(() => {
        loadTasks();
    }, [sorting]);

    const handleCreate = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const deadline = form.deadline ? new Date(form.deadline).toISOString() : undefined;
            await createTask({
                name: form.name,
                description: form.description || undefined,
                deadline: deadline,
                priority: form.priority || undefined,
            });
            setForm({ name: "", description: "", deadline: "", priority: "" });
            await loadTasks();
        } catch (err: any) {
            alert(err.message);
        }
    };


    const handleDelete = async (id: string) => {
        try {
            await deleteTask(id);
            await loadTasks();
        } catch (err: any) {
            alert(err.message);
        }
    };

    const handleToggleStatus = async (taskId: string, currentStatus: task["status"]) => {
        try {
            const newStatus = currentStatus !== "Completed" && currentStatus !== "Late";
            const updatedTask = await toggleTaskStatus(taskId, newStatus);
            setTasks((prevTasks) =>
                prevTasks.map((task) =>
                    task.id === taskId ? updatedTask : task
                )
            );
        } catch (error) {
            console.error("Error toggling task status:", error);
        }
    };

    const handleEdit = (id: string) => {
        const taskToEdit = tasks.find(t => t.id === id);
        if (taskToEdit) {
            setEditingTask(taskToEdit);
            setEditForm({
                name: taskToEdit.name,
                description: taskToEdit.description || "",
                deadline: taskToEdit.deadline ? taskToEdit.deadline.toISOString().split('T')[0] : "",
                priority: taskToEdit.priority || "",
            });
        }
    };

    const handleUpdate = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!editingTask) return;

        try {
            const deadline = editForm.deadline ? new Date(editForm.deadline).toISOString() : undefined;
            await updateTask(editingTask.id, {
                name: editForm.name,
                description: editForm.description || undefined,
                deadline: deadline,
                priority: editForm.priority as any || undefined,
            });
            setEditingTask(null);
            await loadTasks();
        } catch (err: any) {
            alert(err.message);
        }
    };

    return (
        <div className="todolist">
            <form className="get-form" onSubmit={(e) => { e.preventDefault(); loadTasks(); }}>
                <h1>Get All Tasks</h1>
                <select value={sorting} onChange={(e) => setSorting(e.target.value)}>
                    <option value="" disabled hidden>Choose sorting</option>
                    <option value="CreateAsc">Create Time Ascending</option>
                    <option value="CreateDesc">Create Time Descending</option>
                    <option value="PriorityAsc">Priority Ascending</option>
                    <option value="PriorityDesc">Priority Descending</option>
                    <option value="DeadlineAsc">Deadline Ascending</option>
                    <option value="DeadlineDesc">Deadline Descending</option>
                </select>
                <button>Get</button>
            </form>

            <form className="create-task-form" onSubmit={handleCreate}>
                <h1>Create Task</h1>
                <div>
                    <label htmlFor="task-name">Name:</label>
                    <input
                        id="task-name"
                        value={form.name}
                        onChange={e => setForm(f => ({...f, name: e.target.value}))}
                        required
                        minLength={4}
                        aria-invalid={form.name.length > 0 && form.name.length < 4}
                    />
                </div>
                <div>
                    <label htmlFor="task-description">Description:</label>
                    <input
                        id="task-description"
                        value={form.description} 
                        onChange={e => setForm(f => ({ ...f, description: e.target.value }))} 
                    />
                </div>
                <div>
                    <label htmlFor="task-deadline">Deadline:</label>
                    <input 
                        id="task-deadline"
                        type="date" 
                        value={form.deadline} 
                        onChange={e => setForm(f => ({ ...f, deadline: e.target.value }))} 
                    />
                </div>
                <div>
                    <label htmlFor="task-priority">Priority</label>
                    <select 
                        id="task-priority"
                        value={form.priority} 
                        onChange={e => setForm(f => ({ ...f, priority: e.target.value }))}
                    >
                        <option value="" disabled hidden>Choose priority</option>
                        <option value="Low">Low</option>
                        <option value="Medium">Medium</option>
                        <option value="High">High</option>
                        <option value="Critical">Critical</option>
                    </select>
                </div>
                <button>Create</button>
            </form>

            {editingTask && (
                <div className="edit-task-modal">
                    <form className="edit-task-form" onSubmit={handleUpdate}>
                        <h1>Edit Task</h1>
                        <div>
                            <p>Name:</p>
                            <input value={editForm.name} onChange={e => setEditForm(f => ({ ...f, name: e.target.value }))} required />
                        </div>
                        <div>
                            <p>Description:</p>
                            <input value={editForm.description} onChange={e => setEditForm(f => ({ ...f, description: e.target.value }))} />
                        </div>
                        <div>
                            <p>Deadline:</p>
                            <input type="date" value={editForm.deadline} onChange={e => setEditForm(f => ({ ...f, deadline: e.target.value }))} />
                        </div>
                        <div>
                            <p>Priority</p>
                            <select value={editForm.priority} onChange={e => setEditForm(f => ({ ...f, priority: e.target.value }))}>
                                <option value="" disabled hidden>Choose priority</option>
                                <option value="Low">Low</option>
                                <option value="Medium">Medium</option>
                                <option value="High">High</option>
                                <option value="Critical">Critical</option>
                            </select>
                        </div>
                        <div className="edit-form-buttons">
                            <button type="submit">Save</button>
                            <button type="button" onClick={() => setEditingTask(null)}>Cancel</button>
                        </div>
                    </form>
                </div>
            )}

            {tasks.map(task => (
                <Task
                    key={task.id}
                    id={task.id}
                    title={task.name}
                    description={task.description}
                    deadline={task.deadline}
                    priority={task.priority}
                    status={task.status}
                    isDone={task.status === "Completed" || task.status === "Late"}
                    createdAt={task.createdAt}
                    changedAt={task.changedAt}
                    onDelete={handleDelete}
                    onToggleStatus={handleToggleStatus}
                    onEdit={handleEdit}
                />
            ))}
        </div>
    );
};
