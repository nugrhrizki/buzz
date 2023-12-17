import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  createSolidTable,
  getCoreRowModel,
  getExpandedRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
} from "@tanstack/solid-table";
import { Match, ParentProps, Show, Switch, createSignal } from "solid-js";

import { Sender } from "@/models/sender";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import DataTable from "@/components/ui/data-table";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";

type ColumnDefiniton = ColumnDef<Sender>;

export const columns: ColumnDefiniton[] = [
  DataTable.RowExpand(),
  {
    accessorKey: "name",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Nama" />,
    enableHiding: false,
  },
  {
    accessorKey: "token",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Token" />,
  },
  {
    accessorKey: "connected",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Status" />,
    cell: (cell) => (
      <Switch fallback={<Badge variant="secondary">inactive</Badge>}>
        <Match when={cell.row.original.connected === 1}>
          <Badge variant="default">connected</Badge>
        </Match>
        <Match when={cell.row.original.connected === 0}>
          <Badge variant="destructive">disconnected</Badge>
        </Match>
      </Switch>
    ),
  },
];

export function createSenderTable(data: Sender[], columns: ColumnDefiniton[]) {
  const [sorting, setSorting] = createSignal<SortingState>([]);
  const [columnFilters, setColumnFilters] = createSignal<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = createSignal<VisibilityState>({});
  const [rowSelection, setRowSelection] = createSignal({});

  return createSolidTable({
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getRowCanExpand: () => true,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      get sorting() {
        return sorting();
      },
      get columnFilters() {
        return columnFilters();
      },
      get columnVisibility() {
        return columnVisibility();
      },
      get rowSelection() {
        return rowSelection();
      },
    },
  });
}

export function createContactsTable<T>(data: T[], columns: ColumnDef<T>[]) {
  const [sorting, setSorting] = createSignal<SortingState>([]);
  const [columnFilters, setColumnFilters] = createSignal<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = createSignal<VisibilityState>({
    jid: false,
  });
  const [rowSelection, setRowSelection] = createSignal({});

  return createSolidTable({
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getRowCanExpand: () => true,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      get sorting() {
        return sorting();
      },
      get columnFilters() {
        return columnFilters();
      },
      get columnVisibility() {
        return columnVisibility();
      },
      get rowSelection() {
        return rowSelection();
      },
    },
  });
}

interface TableDetailProps {
  sender: Sender;
}

export function TableDetail(props: ParentProps<TableDetailProps>) {
  return (
    <div class="p-4">
      <Card>
        <CardContent class="p-2 pt-2">
          <Table>
            <TableBody>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">ID</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.sender.id || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">JID</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.sender.jid || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Name</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.sender.name || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Token</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.sender.token || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Status</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>
                  <Switch fallback={<Badge variant="secondary">inactive</Badge>}>
                    <Match when={props.sender.connected === 1}>
                      <Badge variant="default">connected</Badge>
                    </Match>
                    <Match when={props.sender.connected === 0}>
                      <Badge variant="destructive">disconnected</Badge>
                    </Match>
                  </Switch>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
        <CardFooter>
          <Show when={props.children}>
            <div class="flex justify-end w-full gap-x-4">{props.children}</div>
          </Show>
        </CardFooter>
      </Card>
    </div>
  );
}
