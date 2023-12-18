import { Table } from "@tanstack/solid-table";
import { RiSystemRefreshLine } from "solid-icons/ri";
import { TbLoader, TbPlus } from "solid-icons/tb";
import { Show, createEffect, createSignal } from "solid-js";

import { useUsers } from "@/services/user";

import { User } from "@/models/user";

import { Button } from "@/components/ui/button";
import Datatable from "@/components/ui/data-table";
import { Input } from "@/components/ui/input";

import { TableDetail, columns, createUserTable } from "./data-table";

function UserPage() {
  const user = useUsers();
  const [table, setTable] = createSignal<Table<User>>(createUserTable(user.data || [], columns));

  createEffect(() => {
    setTable(createUserTable(user.data || [], columns));
  });

  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">
          User
          <Show when={user.isFetching}>
            <TbLoader class="inline-flex ml-2 h-4 w-4 animate-spin" />
          </Show>
        </h2>
      </div>
      <Datatable.Root>
        <div class="flex items-end gap-x-4">
          <Input
            placeholder="Filter username..."
            value={(table().getColumn("username")?.getFilterValue() as string) ?? ""}
            onInput={(event) => {
              table().getColumn("username")?.setFilterValue(event.target.value);
            }}
            class="max-w-sm"
          />
          <div class="ml-auto flex items-end gap-x-4">
            <Button
              variant="outline"
              size="icon"
              onClick={() => {
                user.refetch();
              }}>
              <RiSystemRefreshLine
                classList={{
                  "animate-spin": user.isFetching,
                }}
              />
            </Button>
            <Datatable.ColumnVisibility table={table()} />
            <Button>
              <TbPlus class="w-5 h-5  mr-2" />
              Tambah
            </Button>
          </div>
        </div>
        <Datatable.Table table={table()}>{(row) => <TableDetail user={row.original} />}</Datatable.Table>
        <Datatable.Pagination table={table()} />
      </Datatable.Root>
    </div>
  );
}
export default UserPage;
