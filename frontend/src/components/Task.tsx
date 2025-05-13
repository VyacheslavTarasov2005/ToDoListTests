import type { status } from "../enums/status.ts";
import type { priority } from "../enums/priority.ts";

type Props = {
    title: string;
    description?: string;
    deadline?: Date;
    status: status;
    priority?: priority;
    isDone: boolean;
}

export const Task = ({ title, description, deadline, status, priority, isDone }: Props) => {
    return (
        <div className="task">
            <div className="task-header">
                <h3 className="task-title">{title}</h3>
                {priority && (
                    <span className={`task-priority priority-${priority}`}>
                        {priority}
                    </span>
                )}
            </div>

            {description && (
                <p className="task-description">{description}</p>
            )}

            <div className="task-meta">
                <span className={`task-status status-${status}`}>
                    {status}
                </span>

                {deadline && (
                    <span className="task-deadline">
                        {deadline.toLocaleString()}
                    </span>
                )}
            </div>

            <div className="task-actions">
                <button className={isDone ? "completed" : ""}>
                    {isDone ? "Completed" : "Mark Complete"}
                </button>
                <button>Edit</button>
                <button>Delete</button>
            </div>
        </div>
    )
}