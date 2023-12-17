import { As } from "@kobalte/core";
import { A } from "@solidjs/router";
import { useQueryClient } from "@tanstack/solid-query";
import { Table } from "@tanstack/solid-table";
import { RiSystemRefreshLine } from "solid-icons/ri";
import { TbLoader, TbPlus } from "solid-icons/tb";
import { Show, createEffect, createSignal } from "solid-js";

import { createSenderMutation, deleteSenderMutation, updateSenderMutation, useSender } from "@/services/sender";

import { CreateSender, Sender } from "@/models/sender";

import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button, buttonVariants } from "@/components/ui/button";
import Datatable from "@/components/ui/data-table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";

import { TableDetail, columns, createSenderTable } from "./data-table";
import { CreateSenderForm, EditSenderForm } from "./form";

function SenderPage() {
  const sender = useSender();
  const [table, setTable] = createSignal<Table<Sender>>(createSenderTable(sender.data || [], columns));

  createEffect(() => {
    setTable(createSenderTable(sender.data || [], columns));
  });

  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">
          Sender
          <Show when={sender.isFetching}>
            <TbLoader class="inline-flex ml-2 h-4 w-4 animate-spin" />
          </Show>
        </h2>
      </div>
      <Datatable.Root>
        <div class="flex items-end gap-x-4">
          <Input
            placeholder="Filter name..."
            value={(table().getColumn("name")?.getFilterValue() as string) ?? ""}
            onInput={(event) => {
              table().getColumn("name")?.setFilterValue(event.target.value);
            }}
            class="max-w-sm"
          />
          <div class="ml-auto flex items-end gap-x-4">
            <Button
              variant="outline"
              size="icon"
              onClick={() => {
                sender.refetch();
              }}>
              <RiSystemRefreshLine
                classList={{
                  "animate-spin": sender.isFetching,
                }}
              />
            </Button>
            <Datatable.ColumnVisibility table={table()} />
            <DialogCreateSenderForm />
          </div>
        </div>
        <Datatable.Table table={table()}>
          {(row) => (
            <TableDetail sender={row.original}>
              <A class={buttonVariants()} href={`/sender/${row.original.id}`}>
                Open
              </A>
              <DialogEditSenderForm sender={row.original} />
              <DialogDeleteSenderForm sender={row.original} />
            </TableDetail>
          )}
        </Datatable.Table>
        <Datatable.Pagination table={table()} />
      </Datatable.Root>
    </div>
  );
}

function DialogCreateSenderForm() {
  const queryClient = useQueryClient();
  const createSender = createSenderMutation();

  const [open, setOpen] = createSignal(false);

  function handleSubmit(sender: CreateSender) {
    createSender.mutate(sender, {
      onSuccess: () => {
        setOpen(false);
        queryClient.invalidateQueries({
          queryKey: ["sender"],
        });
      },
    });
  }

  return (
    <Dialog open={open()} onOpenChange={setOpen}>
      <DialogTrigger as={Button}>
        <TbPlus class="w-5 h-5  mr-2" />
        Create
      </DialogTrigger>
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create Sender</DialogTitle>
          <DialogDescription>Enter sender details below.</DialogDescription>
        </DialogHeader>
        <CreateSenderForm onSubmit={handleSubmit} />
      </DialogContent>
    </Dialog>
  );
}

interface DialogEditSenderFormProps {
  sender: Sender;
}

function DialogEditSenderForm(props: DialogEditSenderFormProps) {
  const queryClient = useQueryClient();
  const editSender = updateSenderMutation();

  const [open, setOpen] = createSignal(false);

  function handleSubmit(sender: CreateSender) {
    let payload = {
      id: props.sender.id,
      ...sender,
    };
    editSender.mutate(payload, {
      onSuccess: () => {
        setOpen(false);
        queryClient.invalidateQueries({
          queryKey: ["sender"],
        });
      },
    });
  }

  return (
    <Dialog open={open()} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <As component={Button} variant="outline">
          Edit
        </As>
      </DialogTrigger>
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit Sender</DialogTitle>
          <DialogDescription>Enter sender details below.</DialogDescription>
        </DialogHeader>
        <EditSenderForm onSubmit={handleSubmit} sender={props.sender} />
      </DialogContent>
    </Dialog>
  );
}

interface DialogDeleteSenderFormProps {
  sender: Sender;
}

function DialogDeleteSenderForm(props: DialogDeleteSenderFormProps) {
  const queryClient = useQueryClient();
  const deleteSender = deleteSenderMutation();

  const [open, setOpen] = createSignal(false);

  function handleDelete() {
    deleteSender.mutate(props.sender.id || 0, {
      onSuccess: () => {
        setOpen(false);
        queryClient.invalidateQueries({
          queryKey: ["sender"],
        });
      },
    });
  }

  return (
    <AlertDialog open={open()} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <As component={Button} variant="destructive">
          Delete
        </As>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogTitle>Are you sure?</AlertDialogTitle>
        <AlertDialogDescription>
          You are about to delete <span class="font-medium">"{props.sender.name}"</span> sender. This action cannot be
          undone.
        </AlertDialogDescription>
        <div class="flex justify-end gap-x-2">
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            Delete
          </Button>
        </div>
      </AlertDialogContent>
    </AlertDialog>
  );
}

export default SenderPage;
