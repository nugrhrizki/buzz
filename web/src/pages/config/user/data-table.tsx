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
import { createSignal } from "solid-js";

import { formatDate } from "@/lib/utils";

import { User } from "@/models/user";

import { Card, CardContent } from "@/components/ui/card";
import DataTable from "@/components/ui/data-table";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";

type ColumnDefiniton = ColumnDef<User>;

export const columns: ColumnDefiniton[] = [
  DataTable.RowExpand(),
  DataTable.RowSelect(),
  {
    accessorKey: "name",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Name" />,
    enableHiding: false,
  },
  {
    accessorKey: "username",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Username" />,
  },
  {
    accessorKey: "created_at",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Created At" />,
    cell: (cell) => {
      if (!cell.row.original.created_at) {
        return "-";
      }

      return formatDate(cell.row.original.created_at);
    },
  },
  {
    accessorKey: "updated_at",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Updated At" />,
    cell: (cell) => {
      if (!cell.row.original.created_at) {
        return "-";
      }

      return formatDate(cell.row.original.created_at);
    },
  },
];

export function createUserTable(data: User[], columns: ColumnDefiniton[]) {
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

interface TableDetailProps {
  user: User;
}

export function TableDetail(props: TableDetailProps) {
  return (
    <div class="p-4">
      <Card>
        <CardContent class="p-2 pt-2">
          <Table>
            <TableBody>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">ID</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.user.id || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Name</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.user.name || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Username</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.user.username || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Email</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.user.email || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Whatsapp</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.user.whatsapp || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Created At</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.user.created_at ? formatDate(props.user.created_at) : "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Updated At</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.user.updated_at ? formatDate(props.user.updated_at) : "-"}</TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
