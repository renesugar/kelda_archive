package supervisor

import (
	"net"

	"github.com/kelda/kelda/db"
	"github.com/kelda/kelda/minion/docker"
)

type testCtx struct {
	fd    fakeDocker
	execs [][]string

	conn  db.Conn
	trigg db.Trigger
}

func initTest(r db.Role) *testCtx {
	conn = db.New()
	md, _dk := docker.NewMock()
	ctx := testCtx{fakeDocker{_dk, md}, nil, conn,
		conn.Trigger(db.MinionTable, db.EtcdTable)}
	role = r
	dk = ctx.fd.Client

	ctx.conn.Txn(db.AllTables...).Run(func(view db.Database) error {
		m := view.InsertMinion()
		m.Self = true
		view.Commit(m)
		e := view.InsertEtcd()
		view.Commit(e)
		return nil
	})

	execRun = func(name string, args ...string) error {
		ctx.execs = append(ctx.execs, append([]string{name}, args...))
		return nil
	}

	cfgGateway = func(name string, ip net.IPNet) error {
		execRun("cfgGateway", ip.String())
		return nil
	}

	return &ctx
}

func (ctx *testCtx) run() {
	ctx.execs = nil
	select {
	case <-ctx.trigg.C:
	}

	switch role {
	case db.Master:
		runMasterOnce()
	case db.Worker:
		runWorkerOnce()
	}
}

type fakeDocker struct {
	docker.Client
	md *docker.MockClient
}

func (f fakeDocker) running() map[string][]string {
	containers, _ := f.List(nil)

	res := map[string][]string{}
	for _, c := range containers {
		res[c.Name] = c.Args
	}
	return res
}
