import { Link, useNavigate } from "react-router-dom";
import styles from "./Header.module.css"
import { Container } from "../Container/Container";
import { Button } from "../Button/Button";
import { useState } from "react";

export function Header() {
    const navigate = useNavigate()

    const [isAuthorized, setIsAuthorized] = useState(!!localStorage.getItem("token"))

    const handleLogout = () => {
        localStorage.removeItem("token")

        setIsAuthorized(false)

        navigate("/login")
    }

    const handleAuthorize = () => navigate("/login")

    return (
        <header className={styles.header}>
            <Container className={styles.container}>
                <div>
                    <Link to={isAuthorized ? "/dashboard" : "/login"} className={styles.logotype}>LiteBank</Link>
                </div>

                <nav className={styles.navigation}>
                    {isAuthorized ? (
                        <>
                            <Link to="/dashboard">Счета</Link>
                            <Link to="/transfers">История</Link>
                            <Button variant="danger"
                                        onClick={handleLogout}>
                                            Выйти
                            </Button>
                        </>
                    ) : (
                        <Button variant="primary"
                                onClick={handleAuthorize}>
                                        Войти
                        </Button>
                    )}

                </nav>
            </Container>
        </header>
    )
}