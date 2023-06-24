import React, { useState } from "react";
import { KeyedMutator } from "swr";
import { Todo } from "../types";
import { ENDPOINT } from "../App";
import { Modal, Group, TextInput, Textarea } from "@mantine/core";
import { useForm } from "@mantine/form";

const AddTask = ({ mutate }: { mutate: KeyedMutator<Todo[]> }) => {
	const [open, setOpen] = useState<boolean>(false);

	const form = useForm({
		initialValues: {
			title: "",
			body: ""
		}
	});

	const createTodo = async (values: { title: string, body: string }) => {
		const updated = await fetch(`${ENDPOINT}/api/todos`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify(values)
		}).then((r) => r.json());

		await mutate(updated);
		form.reset();
		setOpen(false);
	};

	return (
		<>
			<Modal opened={open} onClose={() => setOpen(false)} title="Create Task">
				<form onSubmit={form.onSubmit(createTodo)}>
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
						Create Task
					</button>
				</form>
			</Modal>
			<Group position="center">
				<button onClick={() => setOpen(true)} className="p-2 px-4 bg-indigo-500 hover:bg-indigo-600 rounded-md">
					Add Task
				</button>
			</Group>
		</>
	);
};

export default AddTask;
