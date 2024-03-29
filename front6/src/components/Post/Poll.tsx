import * as types from "../../api/api";
import React from "react";
import * as styles from "./Poll.module.css";
import {useQuery, useQueryClient} from "@tanstack/react-query";
import {ApiPolls} from "../../api/client";

export function Poll(props: { pollId: string }) {
    const {data} = useQuery<types.Poll>({
        queryKey: ['poll', props.pollId],
    });
    if (!data) return null;

    const queryClient = useQueryClient();
    const poll = data;

    let isVoted = false;
    for (let answer of poll.answers) {
        isVoted = isVoted || answer.isVoted;
    }

    const onVote = (answerId: string) => {
        if (isVoted) return;

        ApiPolls.Vote({
            pollId: poll.id,
            answerIds: [answerId],
        }).then(() => {
            queryClient.setQueryData(
                ['poll', props.pollId],
                (oldData: types.Poll) => {
                    const copy = structuredClone(oldData) as types.Poll;

                    for (let answer of copy.answers) {
                        if (answer.id == answerId) {
                            answer.isVoted = true;
                            answer.votedCount++;
                        }
                    }

                    queryClient.setQueryData(['poll', props.pollId], copy);
                }
            )
        })
    }

    const onDeleteVote = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        ApiPolls.DeleteVote({
            pollId: poll.id,
        }).then(() => {
            queryClient.setQueryData(
                ['poll', props.pollId],
                (oldData: types.Poll) => {
                    const copy = structuredClone(oldData) as types.Poll;

                    for (let answer of copy.answers) {
                        if (answer.isVoted) {
                            answer.isVoted = false;
                            answer.votedCount--;
                        }
                    }

                    queryClient.setQueryData(['poll', props.pollId], copy);
                }
            )
        })
    }


    return <div className={styles.poll}>
        <div className={styles.question}>
            {poll.question}
        </div>

        <div>
            {poll.answers.map(answer => {
                let className = styles.answer;
                if (answer.isVoted) {
                    className += " " + styles.votedAnswer;
                }

                return <div className={className} key={answer.id} onClick={() => onVote(answer.id)}>
                    {answer.answer}
                    <div className={styles.votersCount}>{answer.votedCount}</div>
                </div>
            })}
            {isVoted && <a href="#" className={styles.deleteVote} onClick={onDeleteVote}>Delete vote</a>}
        </div>
    </div>
}
