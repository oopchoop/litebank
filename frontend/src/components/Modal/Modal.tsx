import React, { useEffect } from "react";
import { Button } from "../Button/Button";
import styles from "./Modal.module.css";

interface ModalProps {
    isOpen: boolean;
    title: string;
    onClose: () => void;
    children: React.ReactNode;
}

export function Modal({ isOpen, title, onClose, children }: ModalProps) {
    useEffect(() => {
        if (isOpen) {
            document.body.style.overflow = 'hidden';
        } else {
            document.body.style.overflow = 'unset';
        }

        return () => {
            document.body.style.overflow = 'unset';
        }
    }, [isOpen])

    if (!isOpen) {
        return null;
    }

    const closeModal = () => {
        onClose();
    }

    return (
        <div className={styles.overlay}>
            <div className={styles.modal}>
                <div className={styles.header}>
                    <h2 className={styles.title}>{title}</h2>
                    <Button variant="danger" onClick={closeModal}>Закрыть</Button>
                </div>

                <div>
                    {children}
                </div>
            </div>
        </div>
    )
}
