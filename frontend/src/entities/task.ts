import type { priority } from "../enums/priority.ts";
import type { status } from "../enums/status.ts";

export class task {
    public id: string;
    public createdAt: Date;
    public changedAt?: Date;
    public name: string;
    public description?: string;
    public deadline?: Date;
    public priority: priority;
    public status: status;
    public isDone: boolean;

    constructor(id: string, name: string, priority: priority, status: status, createdAt: Date, changedAt?: Date, description?: string, deadline?: Date) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.deadline = deadline;
        this.priority = priority;
        this.status = status;
        this.createdAt = createdAt;
        this.changedAt = changedAt;

        this.isDone = status === "Completed" || status === "Late";
    }
}