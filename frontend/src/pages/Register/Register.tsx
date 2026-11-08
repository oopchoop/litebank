import { AuthForm } from "../../components/Forms/AuthForm"
import { LabelField } from "../../components/Field/LabelField"
import { Button } from "../../components/Button/Button"
import { Link } from "react-router-dom"
import api from "../../api/client"
import React, { useState } from "react"
import toast from "react-hot-toast"

export function Register() {
    const [phone, setPhone] = useState('')
    const [firstName, setFirstName] = useState('')
    const [lastName, setLastName] = useState('')
    const [patronymic, setPatronymic] = useState('')
    const [password, setPassword] = useState('')
    const [passwordConfirm, setPasswordConfirm] = useState('')

    const handleRegistration = async (e: React.FormEvent) => {
        e.preventDefault()

        if (password !== passwordConfirm) {
            toast.error('Password not confirmed')            

            return
        }

        try {
            await api.post('/register', {
                  first_name: firstName,
                  last_name: lastName,
                  phone_number: phone,
                  password: password,
                  patronymic: patronymic
            })

            toast.success('Success registration')

            window.location.href = '/login'
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 'Incorrect register data'

            toast.error(errorMessage)            
        }
    }

    return (
        <AuthForm method="POST"
                  header="Регистрация"
                  onSubmit={handleRegistration}>
            <LabelField label="Номер телефона"
                        type="text"
                        value={phone}
                        onChange={(e) => setPhone(e.target.value)}
                        labelCentered={true} />
            <LabelField label="Имя"
                        type="text"
                        value={firstName}
                        onChange={(e) => setFirstName(e.target.value)}
                        labelCentered={true} />
            <LabelField label="Фамилия"
                        type="text"
                        value={lastName}
                        onChange={(e) => setLastName(e.target.value)}
                        labelCentered={true} />
            <LabelField label="Отчество"
                        type="text"
                        value={patronymic}
                        onChange={(e) => setPatronymic(e.target.value)}
                        labelCentered={true} />
            <LabelField label="Пароль"
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        labelCentered={true} />
            <LabelField label="Подтверждение пароля"
                        type="password"
                        value={passwordConfirm}
                        onChange={(e) => setPasswordConfirm(e.target.value)}
                        labelCentered={true} />
            <Button>
                Зарегистрироваться
            </Button>
            <span>
                Уже есть аккаут? <Link className="blue-links" to="/login">Авторизация</Link>
            </span>
        </AuthForm>
    )
}