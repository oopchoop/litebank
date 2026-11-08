import api from "../../api/client";
import type { Account } from "../../api/types";
import React, { useState, useEffect } from "react";
import { Button } from "../../components/Button/Button";
import { AccountCard } from "../../components/AccountCard/AccountCard";
import { Modal } from "../../components/Modal/Modal";
import { LabelField } from "../../components/Field/LabelField";
import { AuthForm } from "../../components/Forms/AuthForm";
import styles from "./Dashboard.module.css";

export function Dashboard() {
    const [accounts, setAccounts] = useState<Account[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [isTransferOpen, setIsTransferOpen] = useState(false);
    const [isDepositOpen, setIsDepositOpen] = useState(false);
    const [selectedAccount, setSelectedAccount] = useState<number | null>(null);
    const [depositSum, setDepositSum] = useState('');
    const [transferSum, setTransferSum] = useState('');
    const [transferAccout, setTransferAccount] = useState('');

    const fetchAccounts = async () => {
        setIsLoading(true)

        try {
            const response = await api.get<Account[]>('/accounts')

            setAccounts(response.data)
        } catch (err: any) {
            console.log(err);
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

            location.reload()            
        } catch (err: any) {
            console.log(err);
        }
    }

    if (isLoading) {
        return <div>Accounts loading...</div>
    }

    const onDepositClick = (accountId: number) => {
        setIsDepositOpen(true);
        setSelectedAccount(accountId);
    }

    const onTransferClick = (accountId: number) => {
        setIsTransferOpen(true);
        setSelectedAccount(accountId);
    }
    
    const handlerDeposit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            await api.post('/deposit', {
                account_to_id: selectedAccount,
                amount: Number(depositSum) * 100,
            });

            setIsDepositOpen(false);

            setDepositSum('');

            fetchAccounts();
        } catch (err: any) {
            console.log(err);
        } 
    }

    const handlerTransfer = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            await api.post('/transfer', {
                account_to_id: Number(transferAccout),
                account_from_id: selectedAccount,
                amount: Number(transferSum) * 100,
            });

            setIsTransferOpen(false);

            setTransferAccount('');
            setTransferSum('');

            fetchAccounts();
        } catch (err: any) {
            console.log(err);
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