import { Header } from "../Header/Header"
import { Outlet } from "react-router-dom"
import { Container } from "../Container/Container"

export function Layout() {
    return (
        <div>
            <Header />
            <main>
                <Container>
                    <Outlet /> 
                </Container>
            </main>
        </div>
    )
}