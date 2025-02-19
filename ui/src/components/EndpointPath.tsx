import miscStyles from "./style/misc.module.sass";
import React from "react";
import styles from './style/EndpointPath.module.sass';

interface EndpointPathProps {
    method: string,
    path: string
}

export const EndpointPath: React.FC<EndpointPathProps> = ({method, path}) => {
    return <div className={styles.container}>
        {method && <span className={`${miscStyles.protocol} ${miscStyles.method}`}>{method}</span>}
        {path && <div title={path} className={styles.path}>{path}</div>}
    </div>
};