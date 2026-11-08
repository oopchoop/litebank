import React from "react";
import styles from "./LabelField.module.css"

interface FieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
    label: string;
    labelCentered?: boolean;
    error?: string;
}

export const LabelField: React.FC<FieldProps> = ({labelCentered = false, label, error, ...props}) => {
    return (
        <div className={styles.container}>
            <label className={labelCentered ? styles.centered : ''}>{label}</label>
            <input  className={styles.field}
                    {...props} />
            {error && <span>{error}</span>}
        </div>
    )
}