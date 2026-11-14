import { Routes, Route, Navigate } from "react-router-dom";
import { Login } from "./pages/Login/Login";
import { Register } from "./pages/Register/Register";
import { Layout } from "./components/Layout/Layout";
import { Dashboard } from "./pages/Dashboard/Dashboard";
import { History } from "./pages/History/History";
import { Toaster } from "react-hot-toast";

function App() {
    return (
        <>
            <Toaster position="bottom-right" />
            <Routes>
                <Route element={<Layout />}>
                    <Route path="/login" element={<Login />}/>
                    <Route path="/register" element={<Register />}/>

                    <Route path="/dashboard" element={<Dashboard />}/>

                    <Route path="/transfers" element={<History />}/>

                    <Route path="*" element={<Navigate to="/login" replace />}/>
                </Route>
            </Routes>
        </>
    )
}

export default App