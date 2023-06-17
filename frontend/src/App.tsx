import React from "react";
import useSWR from "swr";
import { AddTodo } from "./components";
import { Todo } from "./types";
import { CheckCircleFillIcon, CircleIcon } from "@primer/octicons-react";

export const ENDPOINT = "http://localhost:4000";

const fetcher = (url: string) =>
	fetch(`${ENDPOINT}/${url}`).then((r) => r.json());

const App = () => {
	const { data, mutate } = useSWR<Todo[]>("api/todos", fetcher);

	console.log(data);

	const markTodoAddDone = async (id: number) => {
		const updated = await fetch(`${ENDPOINT}/api/todos/${id}/done`, {
			method: "PATCH"
		}).then((r) => r.json());

		await mutate(updated);
	};

	return (
		<>
			<h1 className="my-4 text-center text-3xl font-bold">Todo</h1>
			<div className="my-4 grid place-content-center">
				<ul className="w-1/3 mx-auto flex flex-col gap-4">
					{data?.map((t: Todo) => {
						return (
							<li key={t.id} onClick={() => markTodoAddDone(t.id)} className="flex flex-row gap-2 items-center">
								{t.done ?
									<CheckCircleFillIcon className="w-4 h-4 fill-green-500 mr-2" />
									:
									<CircleIcon className="w-4 h-4 fill-neutral-500 mr-2" />
								}
								<div className="flex flex-col">
									{t.title}
									<span className="text-sm text-neutral-500">
										{t.body}
									</span>
								</div>
							</li>
						);
					})}
				</ul>
			</div>

			<AddTodo mutate={mutate} />
		</>
	);
};

export default App;
