import type {sorting} from "../enums/sorting.ts";
import type {task} from "../entities/task.ts";
import type {priority} from "../enums/priority.ts";

const API_BASE = "http://localhost:8080";

export async function fetchTasks(sorting?: sorting) {
    const url = new URL(`${API_BASE}/tasks`);
    if (sorting) url.searchParams.set("sorting", sorting);

    const response = await fetch(url.toString());
    if (!response.ok) throw new Error("Failed to fetch tasks");

    return response.json();
}

export async function createTask(data: {
    name: string;
    description?: string;
    deadline?: string;
    priority?: string;
}) {
    const response = await fetch(`${API_BASE}/tasks`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const err = await response.json();
        throw new Error(err?.Errors?.message ?? "Failed to create task");
    }

    return response.json();
}

export async function deleteTask(id: string) {
    const response = await fetch(`${API_BASE}/tasks/${id}`, {
        method: "DELETE",
    });

    if (!response.ok) {
        const err = await response.json();
        throw new Error(err?.Errors?.message ?? "Failed to delete task");
    }
}

export async function updateTask(
    id: string,
    data: {
        name: string;
        description?: string;
        deadline?: string | null;
        priority?: priority;
    }
): Promise<task> {
    const response = await fetch(`${API_BASE}/tasks/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        throw new Error("Failed to update task");
    }

    return await response.json();
}

export async function toggleTaskStatus(id: string, isDone: boolean): Promise<task> {
    const response = await fetch(`${API_BASE}/tasks/${id}/toggle`, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ isDone }),
    });

    if (!response.ok) {
        throw new Error("Failed to toggle task status");
    }

    return await response.json();
}
