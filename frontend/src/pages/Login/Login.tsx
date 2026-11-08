import { AuthForm } from "../../components/Forms/AuthForm"
import { LabelField } from "../../components/Field/LabelField"
import { Button } from "../../components/Button/Button"
import { Link } from "react-router-dom"
import api from "../../api/client"
import React, { useState } from "react"
import toast from "react-hot-toast"

export function Login() {
    const [phone, setPhone] = useState('')
    const [password, setPassword] = useState('')

    const handleAuth = async  (e: React.FormEvent) => {
        e.preventDefault()

        try {
            const response = await api.post('/auth', {
                phone_number: phone,
                password: password
            })

            localStorage.setItem('token', response.data.user_token)

            window.location.href = '/dashboard'
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 'Неверный номер телефона или пароль';

            toast.error(errorMessage)
        }
    }

    return (
        <AuthForm method="POST"
                  header="Авторизация"
                  onSubmit={handleAuth}>
            <LabelField label="Номер телефона"
                        type="text"
                        value={phone}
                        onChange={(e) => setPhone(e.target.value)}
                        labelCentered={true} />
            <LabelField label="Пароль"
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        labelCentered={true} />
            <Button>
                Авторизоваться
            </Button>
            <span>
                Еще нет аккаунта? <Link className="blue-links" to="/register">Зарегистрироваться</Link>
            </span>
        </AuthForm>
    )
}