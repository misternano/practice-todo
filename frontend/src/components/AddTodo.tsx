import React, { useState } from "react";
import { KeyedMutator } from "swr";
import { Todo } from "../types";
import { ENDPOINT } from "../App";
import { Modal, Group, TextInput, Textarea } from "@mantine/core";
import { useForm } from "@mantine/form";

const AddTodo = ({ mutate }: { mutate: KeyedMutator<Todo[]> }) => {
	const [open, setOpen] = useState<boolean>(false);

	const form = useForm({
		initialValues: {
			title: "",
			body: "",
		}
	})

	const createTodo = async (values: { title: string, body: string }) => {
		const updated = await fetch(`${ENDPOINT}/api/todos`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(values),
		}).then((r) => r.json());

		mutate(updated);
		form.reset();
		setOpen(false);
	}

	return (
		<>
			<Modal opened={open} onClose={() => setOpen(false)} title="Create Todo">
				<form onSubmit={form.onSubmit(createTodo)}>
					<TextInput
						required
						mb={12}
						label="Todo"
						placeholder="What do you want to do?"
						{...form.getInputProps("title")}
					/>
					<Textarea
						required
						mb={12}
						label="Body"
						placeholder="Tell me more..."
						{...form.getInputProps("body")}
					/>

					<button type="submit">Create Todo</button>
				</form>
			</Modal>

			<Group position="center">
				<button onClick={() => setOpen(true)} className="p-2 px-4 bg-blue-500 hover:bg-blue-600 rounded-md">
					ADD TODO
				</button>
			</Group>
		</>
	)
}

export default AddTodo;