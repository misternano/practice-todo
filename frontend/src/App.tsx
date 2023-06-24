import React, { useState } from "react";
import useSWR from "swr";
import { AddTask } from "./components";
import { Todo } from "./types";
import { CheckCircleFillIcon, CircleIcon, PencilIcon, TrashIcon } from "@primer/octicons-react";
import { Modal, Textarea, TextInput } from "@mantine/core";
import { useForm } from "@mantine/form";

export const ENDPOINT = "http://localhost:4000";

const fetcher = (url: string) =>
	fetch(`${ENDPOINT}/${url}`).then((r) => r.json());

const App = () => {
	const { data, mutate } = useSWR<Todo[]>("api/todos", fetcher);
	const [showDeleteModal, setShowDeleteModal] = useState<boolean>(false);
	const [showEditModal, setShowEditModal] = useState<boolean>(false);
	const [actionItemId, setActionItemId] = useState<number | null>(null);

	const form = useForm({
		initialValues: {
			title: "",
			body: ""
		}
	});

	const toggleTaskDoneState = async (id: number) => {
		const updated = await fetch(`${ENDPOINT}/api/todos/${id}/done`, {
			method: "PATCH"
		}).then((r) => r.json());

		await mutate(updated);
	};

	const handleEditTask = (id: number) => {
		setShowEditModal(true);
		setActionItemId(id);
	};

	const editTask = async (updatedTodo: { title: string, body: string }) => {
		try {
			const edited = await fetch(`${ENDPOINT}/api/todos/${actionItemId}/edit`, {
				method: "PATCH",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify(updatedTodo)
			}).then((r) => r.json());

			await mutate(edited);
			setShowEditModal(false);
			setActionItemId(null);
		} catch (err) {
			console.error(err);
		}
	};

	const handleDeleteConfirmation = (id: number) => {
		setShowDeleteModal(true);
		setActionItemId(id);
	};

	const handleDelete = async () => {
		try {
			const deleted = await fetch(`${ENDPOINT}/api/todos/${actionItemId}/del`, {
				method: "DELETE"
			}).then((r) => r.json());

			await mutate(deleted);
			setShowDeleteModal(false);
			setActionItemId(null);
		} catch (err) {
			console.error(err);
		}
	};

	return (
		<>
			<h1 className="my-4 text-center text-3xl font-bold">Todo</h1>
			<div className="lg:max-w-[50vw] max-w-[90vw] mx-auto my-4 grid place-content-center">
				<ul className={`${data && data.length > 4 ? "grid lg:grid-cols-2" : "flex flex-col"} gap-4`}>
					<Modal opened={showDeleteModal} onClose={() => setShowDeleteModal(false)} title="Delete Confirmation">
						<div className="text-center mb-6 p-2 bg-red-300/50 rounded text-red-900">
							<p className="">Are you sure you want to delete task {actionItemId}?</p>
						</div>
						<div className="flex flex-row gap-2 justify-end">
							<button onClick={() => setShowDeleteModal(false)} className="p-1 px-4 bg-neutral-500 hover:bg-neutral-600 rounded-md text-white">
								Cancel
							</button>
							<button onClick={handleDelete} className="p-1 px-4 bg-red-500 hover:bg-red-700 rounded-md text-white">
								Delete
							</button>
						</div>
					</Modal>
					<Modal opened={showEditModal} onClose={() => setShowEditModal(false)} title={`Edit Task ${actionItemId}`}>
						<form onSubmit={() => editTask(form.values)}>
							<TextInput
								required
								data-autofocus
								mb={16}
								label="Title"
								placeholder="What do you want to do?"
								{...form.getInputProps("title")}
							/>
							<Textarea
								mb={16}
								label="Body"
								placeholder="Tell me more..."
								{...form.getInputProps("body")}
							/>
							<button className="w-full p-1 px-2 bg-indigo-500 hover:bg-indigo-600 rounded text-white" type="submit">
								Update Task
							</button>
						</form>
					</Modal>
					{data && data.length > 0 ? data.map((t: Todo) => {
						return (
							<li className={`group relative ${t.done ? "border-green-600" : "border-neutral-700"} border bg-neutral-800 rounded-lg`} key={t.id}>
								<div className="absolute -top-2.5 right-1 hidden group-hover:flex flex-row items-center bg-neutral-800 border border-neutral-700 rounded-lg">
									<button className="flex items-center" onClick={() => handleEditTask(t.id)}>
										<PencilIcon className="w-6 h-6 p-1 fill-neutral-400 hover:fill-neutral-300 hover:bg-neutral-900/50 rounded-l-lg" />
									</button>
									<button className="flex items-center" onClick={() => handleDeleteConfirmation(t.id)}>
										<TrashIcon className="w-6 h-6 p-1 fill-red-400 hover:fill-red-500 hover:bg-neutral-900/50 rounded-r-lg" />
									</button>
								</div>
								<div className="flex flex-row p-2 gap-2 items-center">
									<div className="w-full flex flex-col gap-2">
										<div className="flex flex-row gap-2">
											<button className="flex items-center" onClick={() => toggleTaskDoneState(t.id)}>
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
										<div className={`${t.body ? "p-2 bg-neutral-700/50 rounded-md select-text" : "hidden"}`}>
											<p className="text-sm text-neutral-300 line-clamp-2">
												{t.body}
											</p>
										</div>
									</div>
								</div>
							</li>
						);
					}) : <p>{data ? "No tasks yet!" : "Could not find any tasks."}</p>}
				</ul>
			</div>

			<AddTask mutate={mutate} />
		</>
	);
};

export default App;
