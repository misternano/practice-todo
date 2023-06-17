import React, { useState } from "react";
import useSWR from "swr";
import { AddTodo } from "./components";
import { Todo } from "./types";
import { CheckCircleFillIcon, CircleIcon, XIcon } from "@primer/octicons-react";
import { Modal } from "@mantine/core";

export const ENDPOINT = "http://localhost:4000";

const fetcher = (url: string) =>
	fetch(`${ENDPOINT}/${url}`).then((r) => r.json());

const App = () => {
	const { data, mutate } = useSWR<Todo[]>("api/todos", fetcher);
	const [showModal, setShowModal] = useState<boolean>(false);
	const [deleteItemId, setDeleteItemId] = useState<number | null>(null);

	const toggleTodoDoneState = async (id: number) => {
		const updated = await fetch(`${ENDPOINT}/api/todos/${id}/done`, {
			method: "PATCH"
		}).then((r) => r.json());

		await mutate(updated);
	};

	const handleDeleteConfirmation = (id: number) => {
		setShowModal(true);
		setDeleteItemId(id);
	};

	const handleDelete = async () => {
		const deleted = await fetch(`${ENDPOINT}/api/todos/${deleteItemId}/del`, {
			method: "DELETE"
		}).then((r) => r.json());

		await mutate(deleted);
		setShowModal(false);
		setDeleteItemId(null);
	};

	return (
		<>
			<h1 className="my-4 text-center text-3xl font-bold">Todo</h1>
			<div className="max-w-[50vw] mx-auto my-4 grid place-content-center">
				<ul className={`${data && data.length < 4 ? "flex flex-col" : "grid grid-cols-2"} gap-4`}>
					<Modal opened={showModal} onClose={() => setShowModal(false)} title="Delete Confirmation">
						<p className="text-center">This cannot be undone!</p>
						<div className="flex flex-row gap-2 justify-center">
							<button onClick={() => setShowModal(false)} className="p-2 px-4 bg-gray-500 hover:bg-gray-600 rounded-md text-white">
								Cancel
							</button>
							<button onClick={handleDelete} className="p-2 px-4 bg-rose-500 hover:bg-rose-600 rounded-md text-white">
								Delete
							</button>
						</div>
					</Modal>
					{data && data.length > 0 ? data.map((t: Todo) => {
						return (
							<li className={`relative ${t.done ? "border-green-600" : "border-neutral-700"} border bg-neutral-800 rounded-lg`} key={t.id}>
								<div className="flex flex-row p-2 gap-2 items-center">
									<button onClick={() => handleDeleteConfirmation(t.id)} className="absolute top-1 right-1">
										<XIcon className="w-5 h-5 p-0.5 hover:fill-red-500 rounded" />
									</button>
									<div className="w-full flex flex-col gap-2">
										<div className="flex flex-row gap-2 items-center">
											<button className="flex items-center" onClick={() => toggleTodoDoneState(t.id)}>
												{t.done ?
													<CheckCircleFillIcon className="w-4 h-4 fill-green-500" />
													:
													<CircleIcon className="w-4 h-4 fill-neutral-500" />
												}
											</button>
											<h2 className="flex flex-row gap-2 text-lg font-medium">
												{t.title}
											</h2>
											<span className="rotate-3 text-sm text-indigo-500">{"#" + t.id}</span>
										</div>
										<p className={`${t.body ? "p-2 bg-neutral-700/50 rounded-md text-sm text-neutral-300 select-all" : "hidden"}`}>
											{t.body}
										</p>
									</div>
								</div>
							</li>
						);
					}) : <p>{data ? "No tasks yet!" : "Could not find tasks"}</p>}
				</ul>
			</div>

			<AddTodo mutate={mutate} />
		</>
	);
};

export default App;
