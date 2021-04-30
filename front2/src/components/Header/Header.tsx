import React from "react";
import styles from "./Header.module.css";
import {Link} from "../Link/Link";
import {FeedGetHeaderRequest, FeedGetHeaderResponse, HeaderRenderer} from "../../api/api2";
import {api} from "../../Api";

interface State {
    data: HeaderRenderer;
}

export class Header extends React.Component {
    state: State;

    constructor(props: any) {
        super(props);
        this.state = {
            data: {
                isAuthorized: false,
                mainUrl: "/",
                userName: "Макс",
                userAvatar: "https://sun2.43222.userapi.com/s/v1/ig2/WVgdYwZ6Cd8mMcMunD_Po2YDv0_2BHRGf3ofZ1NHyGcbd9nKDQJ029FOYIgwo614Rqv3RT1hO7z5t01DUSRaJosq.jpg?size=100x0&quality=96&crop=122,105,561,561&ava=1",
            }
        };
    }

    componentDidMount() {
        this.refreshData();
        setInterval(this.refreshData, 60 * 1000);
    }

    refreshData = () => {
        api<FeedGetHeaderRequest, FeedGetHeaderResponse>("meme.Feed", "GetHeader", {}).then(r => {
            this.setState({data: r.renderer});
        })
    }

    render() {
        const data = this.state.data;

        return (
            <div className={styles.Header}>
                <Link href={data.mainUrl} className={styles.Logo}>meme</Link>

                <div className={styles.RightContainer}>
                    {data.isAuthorized && <div className={styles.Name}>{data.userName}</div>}
                    {data.isAuthorized && <img className={styles.Avatar} alt="" src={data.userAvatar}/>}

                    {!data.isAuthorized && <Link className={styles.Name} href={"/login"}>Войти</Link>}
                </div>
            </div>
        );
    }
}
