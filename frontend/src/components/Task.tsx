import type {status} from "../enums/status.ts";
import type {priority} from "../enums/priority.ts";

type Props = {
    id: string;
    title: string;
    description?: string;
    deadline?: Date;
    status: status;
    priority?: priority;
    isDone: boolean;
    createdAt: Date;
    changedAt?: Date;
    onDelete: (id: string) => void;
    onToggleStatus: (id: string, status: status) => void;
    onEdit: (id: string) => void;
};

export const Task = ({ 
    id, 
    title, 
    description, 
    deadline, 
    status, 
    priority, 
    isDone, 
    createdAt,
    changedAt,
    onDelete, 
    onToggleStatus, 
    onEdit 
}: Props) => {
    const getDeadlineStatus = () => {
        if (status === "Overdue") return 'overdue';
        if (!deadline || isDone) return '';
        
        const now = new Date();
        const deadlineDate = new Date(deadline);
        const diffTime = deadlineDate.getTime() - now.getTime();
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        if (diffDays <= 3) return 'urgent';
        return '';
    };

    const formatDate = (date: Date) => {
        return new Date(date).toLocaleString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const deadlineStatus = getDeadlineStatus();

    return (
        <div 
            className={`task ${deadlineStatus}`}
            data-status={status.toLowerCase()}
            data-priority={priority?.toLowerCase()}
            data-deadline={deadline?.toISOString().split('T')[0]}
        >
            <div className="task-header">
                <h3 className="task-title">{title}</h3>
                {priority && (
                    <span className={`task-priority priority-${priority}`}>
                        {priority}
                    </span>
                )}
            </div>

            {description && <p className="task-description">{description}</p>}

            <div className="task-meta">
                <span className={`task-status status-${status}`}>{status}</span>
                {deadline && <span className="task-deadline">Дедлайн: {formatDate(deadline)}</span>}
            </div>

            <div className="task-dates">
                <span className="task-date">Создано: {formatDate(createdAt)}</span>
                {changedAt && <span className="task-date">Изменено: {formatDate(changedAt)}</span>}
            </div>

            <div className="task-actions">
                <button 
                    onClick={() => onToggleStatus(id, status)} 
                    className={isDone ? "completed" : ""}
                    role="checkbox"
                    aria-checked={isDone}
                >
                    {isDone ? "Completed" : "Mark Complete"}
                </button>
                <button onClick={() => onEdit(id)}>Edit</button>
                <button onClick={() => onDelete(id)}>Delete</button>
            </div>
        </div>
    );
};
