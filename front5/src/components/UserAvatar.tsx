import React from "react";
import styles from "./UserAvatar.module.css";
import {Global} from "../store/store";
import {connect} from "react-redux";
import * as types from "../api/types";

export type Props = {
    userId: string;
    width: number;
}

export interface ComponentProps {
    user: types.User;
    online: types.Online;
    width: number;
}

export function Component(props: ComponentProps) {
    const style = {
        width: props.width + 'px',
        height: props.width + 'px',
    };

    let styleOnline = {
        right: '2px',
        bottom: '2px',
    }

    if (props.width == 100) {
        styleOnline = {
            right: '8px',
            bottom: '8px',
        }
    }

    return <div className={styles.container}>
        {props.user.avatar ?
            <img src={props.user.avatar} className={styles.avatar}
                 alt={"Avatar " + props.user.name} style={style}
            /> :
            <div className={styles.avatar} style={style}/>
        }

        {props.online?.isOnline && <div className={styles.online} style={styleOnline}/>}
    </div>;
}

export const UserAvatar = connect((state: Global, ownProps: Props) => {
    return {
        user: state.users.byId[ownProps.userId],
        online: state.online.byId[ownProps.userId],
        width: ownProps.width,
    } as ComponentProps
})(Component);
