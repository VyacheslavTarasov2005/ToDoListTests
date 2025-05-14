import { ToDoList } from "./components/ToDoList.tsx";
import { Header } from "./components/Header.tsx";

function App() {
    return (
        <div className="layout">
            <Header/>
            <div className="page-content">
                <ToDoList />
            </div>
        </div>
    )
}

export default App;