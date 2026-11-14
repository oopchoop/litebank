import api from "../../api/client";
import type { Account } from "../../api/types";
import React, { useState, useEffect } from "react";
import { Button } from "../../components/Button/Button";
import { AccountCard } from "../../components/AccountCard/AccountCard";
import { Modal } from "../../components/Modal/Modal";
import { LabelField } from "../../components/Field/LabelField";
import { AuthForm } from "../../components/Forms/AuthForm";
import styles from "./Dashboard.module.css";
import toast from "react-hot-toast";
import { v4 as uuidv4 } from "uuid";

export function Dashboard() {
    const [accounts, setAccounts] = useState<Account[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isTransferOpen, setIsTransferOpen] = useState(false);
    const [isDepositOpen, setIsDepositOpen] = useState(false);
    const [selectedAccount, setSelectedAccount] = useState<number | null>(null);
    const [depositSum, setDepositSum] = useState('');
    const [transferSum, setTransferSum] = useState('');
    const [transferAccout, setTransferAccount] = useState('');
    const [idempotencyKey, setIdempotencyKey] = useState('');

    const fetchAccounts = async () => {
        setIsLoading(true)

        try {
            const response = await api.get<Account[]>('/accounts')

            setAccounts(response.data)
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 'Fetch accounts error'

            toast.error(errorMessage)
        } finally {
            setIsLoading(false)
        }
    }

    useEffect(() => {
        fetchAccounts()
    }, [])

    const handleCreateAccount = async () => {
        try {
            await api.post('/accounts')

            fetchAccounts()
            
            toast.success('Success account create')
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 'Account create error'

            toast.error(errorMessage)
        }
    }

    if (isLoading) {
        return <div>Accounts loading...</div>
    }

    const onDepositClick = (accountId: number) => {
        setIsDepositOpen(true);
        setSelectedAccount(accountId);
        setIdempotencyKey(uuidv4())
    }

    const onTransferClick = (accountId: number) => {
        setIsTransferOpen(true);
        setSelectedAccount(accountId);
        setIdempotencyKey(uuidv4())
    }
    
    const handlerDeposit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            await api.post('/deposit', {
                account_to_id: selectedAccount,
                amount: Number(depositSum) * 100,
            }, {
                headers: {
                    "Idempotency-Key": idempotencyKey
                },
            });

            setIsDepositOpen(false);

            setDepositSum('');

            toast.success('Success deposit create')

            fetchAccounts();
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 'Account create error'

            toast.error(errorMessage)
        } 
    }

    const handlerTransfer = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            await api.post('/transfer', {
                account_to_id: Number(transferAccout),
                account_from_id: selectedAccount,
                amount: Number(transferSum) * 100,
            }, {
                headers: {
                    "Idempotency-Key": idempotencyKey
                },
            });

            setIsTransferOpen(false);

            setTransferAccount('');
            setTransferSum('');

            toast.success('Success transfer create')

            fetchAccounts();
        } catch (err: any) {
            const errorMessage = err.response?.data?.error || 'Account create error'

            toast.error(errorMessage)
        }
    }

    return (
        <div className={styles.container}>
            <Button onClick={handleCreateAccount}>Создать счет</Button>
            {accounts.length === 0 ? (
                <div>You don`t have accounts</div>
            ) : (
                accounts.map((item) => (
                    <AccountCard key={item.id}
                                 account={item} 
                                 onDepositClick={() => onDepositClick(item.id)} 
                                 onTransferClick={() => onTransferClick(item.id)} 
                    />
                ))
            )}

            <Modal isOpen={isTransferOpen} onClose={() => setIsTransferOpen(false)} title="Перевод">
                <AuthForm onSubmit={handlerTransfer}>
                    <LabelField label="Номер счета"
                                type="number"
                                value={transferAccout} 
                                onChange={(e) => setTransferAccount(e.target.value)}/>
                    <LabelField label="Сумма"
                                type="number"
                                value={transferSum} 
                                onChange={(e) => setTransferSum(e.target.value)}/>
                    <Button type="submit">Перевести</Button>
                </AuthForm>
            </Modal>

            <Modal isOpen={isDepositOpen} onClose={() => setIsDepositOpen(false)} title="Пополнение счета">
                <AuthForm onSubmit={handlerDeposit}>
                    <LabelField label="Сумма"
                                type="number"
                                value={depositSum} 
                                onChange={(e) => setDepositSum(e.target.value)}/>
                    <Button type="submit">Пополнить</Button>
                </AuthForm>
            </Modal>
        </div>
    )
}