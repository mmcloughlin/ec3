package efd

var expectedcosts = []struct {
	ID        string
	Operation string
	Cost      string
}{
	{
		ID:        "g12o/edwards/w-1/doubling/dbl-2008-blr-1",
		Operation: "doubling",
		Cost:      "1I + 2S + 1*d1",
	},
	{
		ID:        "g12o/edwards/w-1/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "1I + 1M + 2S + 1*d2overd1plus1",
	},
	{
		ID:        "g12o/edwards/w-1/diffadd/dadd-2008-blr-1",
		Operation: "diffadd",
		Cost:      "1I + 1M + 2S + 1*d1",
	},
	{
		ID:        "g12o/edwards/w-1/diffadd/dadd-2008-blr",
		Operation: "diffadd",
		Cost:      "1I + 3M + 1S + 1*d2overd1plus1",
	},
	{
		ID:        "g12o/edwards/w-1/ladder/ladd-2008-blr-1",
		Operation: "ladder",
		Cost:      "2I + 1M + 3S + 2*d1",
	},
	{
		ID:        "g12o/edwards/w-1/ladder/ladd-2008-blr",
		Operation: "ladder",
		Cost:      "2I + 4M + 3S + 2*d2overd1plus1",
	},
	{
		ID:        "g12o/edwards/w-1/scaling/copy",
		Operation: "scaling",
		Cost:      "0M",
	},
	{
		ID:        "g12o/edwards/w/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "1I + 1M + 2S + 1*d2overd1plus1",
	},
	{
		ID:        "g12o/edwards/w/diffadd/dadd-2008-blr",
		Operation: "diffadd",
		Cost:      "1I + 3M + 1S + 1*d2overd1plus1",
	},
	{
		ID:        "g12o/edwards/w/ladder/ladd-2008-blr",
		Operation: "ladder",
		Cost:      "2I + 4M + 3S + 2*d2overd1plus1",
	},
	{
		ID:        "g12o/edwards/w/scaling/copy",
		Operation: "scaling",
		Cost:      "0M",
	},
	{
		ID:        "g12o/edwards/wz-1/doubling/dbl-2008-blr-1",
		Operation: "doubling",
		Cost:      "1M + 3S + 1*d1",
	},
	{
		ID:        "g12o/edwards/wz-1/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "1M + 3S + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz-1/diffadd/mdadd-2008-blr-1",
		Operation: "diffadd",
		Cost:      "5M + 1S + 1*d1",
	},
	{
		ID:        "g12o/edwards/wz-1/diffadd/mdadd-2008-blr",
		Operation: "diffadd",
		Cost:      "6M + 1S + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz-1/diffadd/dadd-2008-blr-3",
		Operation: "diffadd",
		Cost:      "6M + 2S + 2*d1",
	},
	{
		ID:        "g12o/edwards/wz-1/diffadd/dadd-2008-blr-2",
		Operation: "diffadd",
		Cost:      "6M + 2S + 1*d1 + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz-1/diffadd/dadd-2008-blr-1",
		Operation: "diffadd",
		Cost:      "7M + 1S + 1*d1",
	},
	{
		ID:        "g12o/edwards/wz-1/diffadd/dadd-2008-blr",
		Operation: "diffadd",
		Cost:      "8M + 1S + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz-1/ladder/mladd-2008-blr-1",
		Operation: "ladder",
		Cost:      "5M + 4S + 2*d1",
	},
	{
		ID:        "g12o/edwards/wz-1/ladder/mladd-2008-blr",
		Operation: "ladder",
		Cost:      "6M + 4S + 1*ee + 1*ff + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz-1/ladder/ladd-2008-blr-1",
		Operation: "ladder",
		Cost:      "7M + 4S + 2*d1",
	},
	{
		ID:        "g12o/edwards/wz-1/ladder/ladd-2008-blr",
		Operation: "ladder",
		Cost:      "8M + 4S + 1*ee + 1*ff + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz-1/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 1M",
	},
	{
		ID:        "g12o/edwards/wz/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "1M + 3S + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz/diffadd/mdadd-2008-blr",
		Operation: "diffadd",
		Cost:      "6M + 1S + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz/diffadd/dadd-2008-blr-2",
		Operation: "diffadd",
		Cost:      "6M + 2S + 1*d1 + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz/diffadd/dadd-2008-blr",
		Operation: "diffadd",
		Cost:      "8M + 1S + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz/ladder/mladd-2008-blr",
		Operation: "ladder",
		Cost:      "6M + 4S + 1*ee + 1*ff + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz/ladder/ladd-2008-blr",
		Operation: "ladder",
		Cost:      "8M + 4S + 1*ee + 1*ff + 1*e + 1*f",
	},
	{
		ID:        "g12o/edwards/wz/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 1M",
	},
	{
		ID:        "g12o/edwards/xy-1/addition/add-2008-blr",
		Operation: "addition",
		Cost:      "2I + 8M + 2S + 2*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xy-1/addition/add-2008-blr",
		Operation: "readdition",
		Cost:      "2I + 7M + 2S + 2*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xy-1/doubling/dbl-2008-blr+1",
		Operation: "doubling",
		Cost:      "1I + 1M + 4S + 2*d1",
	},
	{
		ID:        "g12o/edwards/xy-1/doubling/dbl-2008-blr+cse",
		Operation: "doubling",
		Cost:      "1I + 2M + 4S + 1*d2 + 1*d2d1",
	},
	{
		ID:        "g12o/edwards/xy-1/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "4I + 4M + 10S + 6^4 + 4*d2",
	},
	{
		ID:        "g12o/edwards/xy-1/scaling/copy",
		Operation: "scaling",
		Cost:      "0M",
	},
	{
		ID:        "g12o/edwards/xy/addition/add-2008-blr",
		Operation: "addition",
		Cost:      "2I + 8M + 2S + 2*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xy/addition/add-2008-blr",
		Operation: "readdition",
		Cost:      "2I + 7M + 2S + 2*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xy/doubling/dbl-2008-blr+cse",
		Operation: "doubling",
		Cost:      "1I + 2M + 4S + 1*d2 + 1*d2d1",
	},
	{
		ID:        "g12o/edwards/xy/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "4I + 4M + 10S + 6^4 + 4*d2",
	},
	{
		ID:        "g12o/edwards/xy/scaling/copy",
		Operation: "scaling",
		Cost:      "0M",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/madd-2008-blr",
		Operation: "addition",
		Cost:      "13M + 3S + 2*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/madd-2008-blr",
		Operation: "readdition",
		Cost:      "13M + 1S + 2*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-3",
		Operation: "addition",
		Cost:      "16M + 1S + 1*d1d1 + 3*d1",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-3",
		Operation: "readdition",
		Cost:      "16M + 1S + 1*d1d1 + 3*d1",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-5",
		Operation: "addition",
		Cost:      "16M + 2S + 3*d1",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-5",
		Operation: "readdition",
		Cost:      "16M + 2S + 3*d1",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-2",
		Operation: "addition",
		Cost:      "18M + 2S + 1*d1d1 + 3*d1 + 1*d2 + 2*d2plusd1",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-2",
		Operation: "readdition",
		Cost:      "18M + 2S + 1*d1d1 + 3*d1 + 1*d2 + 2*d2plusd1",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-4",
		Operation: "addition",
		Cost:      "18M + 3S + 3*d1 + 1*d2 + 2*d2plusd1",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-4",
		Operation: "readdition",
		Cost:      "18M + 3S + 3*d1 + 1*d2 + 2*d2plusd1",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-1",
		Operation: "addition",
		Cost:      "21M + 1S + 3*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xyz-1/addition/add-2008-blr-1",
		Operation: "readdition",
		Cost:      "20M + 1S + 2*d1",
	},
	{
		ID:        "g12o/edwards/xyz-1/doubling/dbl-2008-blr-2",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*d1 + 1*sqrtd1",
	},
	{
		ID:        "g12o/edwards/xyz-1/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "2M + 6S + 1*d1 + 1*d2 + 1*d2d1",
	},
	{
		ID:        "g12o/edwards/xyz-1/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g12o/edwards/xyz/addition/madd-2008-blr",
		Operation: "addition",
		Cost:      "13M + 3S + 2*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xyz/addition/madd-2008-blr",
		Operation: "readdition",
		Cost:      "13M + 1S + 2*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xyz/addition/add-2008-blr-2",
		Operation: "addition",
		Cost:      "18M + 2S + 1*d1d1 + 3*d1 + 1*d2 + 2*d2plusd1",
	},
	{
		ID:        "g12o/edwards/xyz/addition/add-2008-blr-2",
		Operation: "readdition",
		Cost:      "18M + 2S + 1*d1d1 + 3*d1 + 1*d2 + 2*d2plusd1",
	},
	{
		ID:        "g12o/edwards/xyz/addition/add-2008-blr-4",
		Operation: "addition",
		Cost:      "18M + 3S + 3*d1 + 1*d2 + 2*d2plusd1",
	},
	{
		ID:        "g12o/edwards/xyz/addition/add-2008-blr-4",
		Operation: "readdition",
		Cost:      "18M + 3S + 3*d1 + 1*d2 + 2*d2plusd1",
	},
	{
		ID:        "g12o/edwards/xyz/addition/add-2008-blr-1",
		Operation: "addition",
		Cost:      "21M + 1S + 3*d1 + 1*d2",
	},
	{
		ID:        "g12o/edwards/xyz/addition/add-2008-blr-1",
		Operation: "readdition",
		Cost:      "20M + 1S + 2*d1",
	},
	{
		ID:        "g12o/edwards/xyz/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "2M + 6S + 1*d1 + 1*d2 + 1*d2d1",
	},
	{
		ID:        "g12o/edwards/xyz/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g12o/hessian/standard/addition/mmadd",
		Operation: "addition",
		Cost:      "8M",
	},
	{
		ID:        "g12o/hessian/standard/addition/mmadd",
		Operation: "readdition",
		Cost:      "7M",
	},
	{
		ID:        "g12o/hessian/standard/addition/madd",
		Operation: "addition",
		Cost:      "10M",
	},
	{
		ID:        "g12o/hessian/standard/addition/madd",
		Operation: "readdition",
		Cost:      "10M",
	},
	{
		ID:        "g12o/hessian/standard/addition/add-2001-jq",
		Operation: "addition",
		Cost:      "12M",
	},
	{
		ID:        "g12o/hessian/standard/addition/add-2001-jq",
		Operation: "readdition",
		Cost:      "12M",
	},
	{
		ID:        "g12o/hessian/standard/addition/add2",
		Operation: "addition",
		Cost:      "12M",
	},
	{
		ID:        "g12o/hessian/standard/addition/add2",
		Operation: "readdition",
		Cost:      "12M",
	},
	{
		ID:        "g12o/hessian/standard/addition/add",
		Operation: "addition",
		Cost:      "12M + 6S",
	},
	{
		ID:        "g12o/hessian/standard/addition/add",
		Operation: "readdition",
		Cost:      "9M + 3S",
	},
	{
		ID:        "g12o/hessian/standard/doubling/dbl2",
		Operation: "doubling",
		Cost:      "6M + 3S",
	},
	{
		ID:        "g12o/hessian/standard/doubling/dbl-2007-hcd-2",
		Operation: "doubling",
		Cost:      "7M + 1S",
	},
	{
		ID:        "g12o/hessian/standard/doubling/dbl-2007-hcd",
		Operation: "doubling",
		Cost:      "7M + 1S",
	},
	{
		ID:        "g12o/hessian/standard/doubling/dbl-2001-jq",
		Operation: "doubling",
		Cost:      "12M",
	},
	{
		ID:        "g12o/hessian/standard/doubling/dbl",
		Operation: "doubling",
		Cost:      "3M + 6^3",
	},
	{
		ID:        "g12o/hessian/standard/tripling/tpl-2007-hcd-2",
		Operation: "tripling",
		Cost:      "11M + 4S + 2*a",
	},
	{
		ID:        "g12o/hessian/standard/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "10M + 1S + 29^3 + 2*d",
	},
	{
		ID:        "g12o/hessian/standard/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g12o/shortw/affine/addition/add",
		Operation: "addition",
		Cost:      "1I + 2M + 1S",
	},
	{
		ID:        "g12o/shortw/affine/addition/add",
		Operation: "readdition",
		Cost:      "1I + 2M + 1S",
	},
	{
		ID:        "g12o/shortw/affine/doubling/dbl",
		Operation: "doubling",
		Cost:      "1I + 2M + 1S",
	},
	{
		ID:        "g12o/shortw/affine/scaling/copy",
		Operation: "scaling",
		Cost:      "0M",
	},
	{
		ID:        "g12o/shortw/extended-0/addition/madd-2005-dl",
		Operation: "addition",
		Cost:      "9M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/extended-0/addition/madd-2005-dl",
		Operation: "readdition",
		Cost:      "9M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/extended-0/addition/add-2005-dl",
		Operation: "addition",
		Cost:      "14M + 3S",
	},
	{
		ID:        "g12o/shortw/extended-0/addition/add-2005-dl",
		Operation: "readdition",
		Cost:      "14M + 3S",
	},
	{
		ID:        "g12o/shortw/extended-0/doubling/mdbl-2008-blr",
		Operation: "doubling",
		Cost:      "2M + 3S + 1*a6",
	},
	{
		ID:        "g12o/shortw/extended-0/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*a6 + 1*sqrta6",
	},
	{
		ID:        "g12o/shortw/extended-1/addition/madd-2007-kk",
		Operation: "addition",
		Cost:      "8M + 4S",
	},
	{
		ID:        "g12o/shortw/extended-1/addition/madd-2007-kk",
		Operation: "readdition",
		Cost:      "8M + 4S",
	},
	{
		ID:        "g12o/shortw/extended-1/addition/madd-2005-dl",
		Operation: "addition",
		Cost:      "8M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/extended-1/addition/madd-2005-dl",
		Operation: "readdition",
		Cost:      "8M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/extended-1/addition/add-2005-dl",
		Operation: "addition",
		Cost:      "13M + 3S",
	},
	{
		ID:        "g12o/shortw/extended-1/addition/add-2005-dl",
		Operation: "readdition",
		Cost:      "13M + 3S",
	},
	{
		ID:        "g12o/shortw/extended-1/doubling/dbl-2008-blr",
		Operation: "doubling",
		Cost:      "2M + 4S + 1*a6 + 1*sqrta6",
	},
	{
		ID:        "g12o/shortw/extended-1/doubling/dbl-2007-kk",
		Operation: "doubling",
		Cost:      "2M + 5S + 2*a6",
	},
	{
		ID:        "g12o/shortw/jacobian/addition/madd-2008-bl",
		Operation: "addition",
		Cost:      "10M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/jacobian/addition/madd-2008-bl",
		Operation: "readdition",
		Cost:      "10M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/jacobian/addition/add-2005-dl",
		Operation: "addition",
		Cost:      "14M + 5S + 1*a2",
	},
	{
		ID:        "g12o/shortw/jacobian/addition/add-2005-dl",
		Operation: "readdition",
		Cost:      "13M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/jacobian/doubling/mdbl-2008-bl",
		Operation: "doubling",
		Cost:      "1M + 2S + 1*a6",
	},
	{
		ID:        "g12o/shortw/jacobian/doubling/dbl-2005-dl",
		Operation: "doubling",
		Cost:      "4M + 5S + 1*a6",
	},
	{
		ID:        "g12o/shortw/jacobian/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 3M + 1S",
	},
	{
		ID:        "g12o/shortw/lambda/addition/add-2013-olar",
		Operation: "addition",
		Cost:      "11M + 2S",
	},
	{
		ID:        "g12o/shortw/lambda/addition/add-2013-olar",
		Operation: "readdition",
		Cost:      "11M + 2S",
	},
	{
		ID:        "g12o/shortw/lambda/doubling/dbl-2013-olar-2",
		Operation: "doubling",
		Cost:      "3M + 4S + 1*a2 + 1*a226 + 1*a21",
	},
	{
		ID:        "g12o/shortw/lambda/doubling/dbl-2013-olar",
		Operation: "doubling",
		Cost:      "4M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/addition/mmadd-2005-dl",
		Operation: "addition",
		Cost:      "5M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/addition/mmadd-2005-dl",
		Operation: "readdition",
		Cost:      "5M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/addition/madd-2005-dl",
		Operation: "addition",
		Cost:      "8M + 5S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/addition/madd-2005-dl",
		Operation: "readdition",
		Cost:      "8M + 5S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/addition/add-2005-dl",
		Operation: "addition",
		Cost:      "13M + 4S",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/addition/add-2005-dl",
		Operation: "readdition",
		Cost:      "13M + 3S",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/doubling/mdbl-2005-dl",
		Operation: "doubling",
		Cost:      "1M + 3S + 1*a2 + 1*a6",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/doubling/dbl-2005-dl-a2-0",
		Operation: "doubling",
		Cost:      "3M + 5S + 1*sqrta6",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/doubling/dbl-2005-dl",
		Operation: "doubling",
		Cost:      "3M + 5S + 1*a2 + 1*a6",
	},
	{
		ID:        "g12o/shortw/lopezdahab-0/doubling/dbl-2005-l",
		Operation: "doubling",
		Cost:      "4M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/addition/mmadd-2005-dl",
		Operation: "addition",
		Cost:      "5M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/addition/mmadd-2005-dl",
		Operation: "readdition",
		Cost:      "5M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/addition/madd-2005-dl",
		Operation: "addition",
		Cost:      "8M + 5S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/addition/madd-2005-dl",
		Operation: "readdition",
		Cost:      "8M + 5S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/addition/add-2005-dl",
		Operation: "addition",
		Cost:      "13M + 4S",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/addition/add-2005-dl",
		Operation: "readdition",
		Cost:      "13M + 3S",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/doubling/mdbl-2005-dl",
		Operation: "doubling",
		Cost:      "1M + 3S + 1*a2 + 1*a6",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/doubling/dbl-2005-dl-a2-1",
		Operation: "doubling",
		Cost:      "3M + 5S + 1*sqrta6",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/doubling/dbl-2005-dl",
		Operation: "doubling",
		Cost:      "3M + 5S + 1*a2 + 1*a6",
	},
	{
		ID:        "g12o/shortw/lopezdahab-1/doubling/dbl-2005-l",
		Operation: "doubling",
		Cost:      "4M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab/addition/mmadd-2005-dl",
		Operation: "addition",
		Cost:      "5M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab/addition/mmadd-2005-dl",
		Operation: "readdition",
		Cost:      "5M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab/addition/madd-2005-dl",
		Operation: "addition",
		Cost:      "8M + 5S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab/addition/madd-2005-dl",
		Operation: "readdition",
		Cost:      "8M + 5S + 1*a2",
	},
	{
		ID:        "g12o/shortw/lopezdahab/addition/add-2005-dl",
		Operation: "addition",
		Cost:      "13M + 4S",
	},
	{
		ID:        "g12o/shortw/lopezdahab/addition/add-2005-dl",
		Operation: "readdition",
		Cost:      "13M + 3S",
	},
	{
		ID:        "g12o/shortw/lopezdahab/doubling/mdbl-2005-dl",
		Operation: "doubling",
		Cost:      "1M + 3S + 1*a2 + 1*a6",
	},
	{
		ID:        "g12o/shortw/lopezdahab/doubling/dbl-2005-dl",
		Operation: "doubling",
		Cost:      "3M + 5S + 1*a2 + 1*a6",
	},
	{
		ID:        "g12o/shortw/lopezdahab/doubling/dbl-2005-l",
		Operation: "doubling",
		Cost:      "4M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/mmadd-2008-bl",
		Operation: "addition",
		Cost:      "7M + 1S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/mmadd-2008-bl",
		Operation: "readdition",
		Cost:      "7M + 1S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/madd-2008-bl",
		Operation: "addition",
		Cost:      "11M + 1S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/madd-2008-bl",
		Operation: "readdition",
		Cost:      "11M + 1S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/madd-2005-dl",
		Operation: "addition",
		Cost:      "11M + 2S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/madd-2005-dl",
		Operation: "readdition",
		Cost:      "11M + 2S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/add-2008-bl",
		Operation: "addition",
		Cost:      "14M + 1S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/add-2008-bl",
		Operation: "readdition",
		Cost:      "14M + 1S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/add-2005-dl-2",
		Operation: "addition",
		Cost:      "15M + 2S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/add-2005-dl-2",
		Operation: "readdition",
		Cost:      "15M + 2S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/add-2005-dl",
		Operation: "addition",
		Cost:      "15M + 2S + 1^3 + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/addition/add-2005-dl",
		Operation: "readdition",
		Cost:      "15M + 2S + 1^3 + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/doubling/mdbl-2005-dl",
		Operation: "doubling",
		Cost:      "5M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/doubling/dbl-2008-bl",
		Operation: "doubling",
		Cost:      "7M + 3S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/doubling/dbl-2005-dl",
		Operation: "doubling",
		Cost:      "7M + 4S + 1*a2",
	},
	{
		ID:        "g12o/shortw/projective/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g12o/shortw/xz/doubling/mdbl-2003-s",
		Operation: "doubling",
		Cost:      "2S",
	},
	{
		ID:        "g12o/shortw/xz/doubling/dbl-2003-s-3",
		Operation: "doubling",
		Cost:      "1M + 3S + 1*sqrta6",
	},
	{
		ID:        "g12o/shortw/xz/doubling/dbl-2003-s-4",
		Operation: "doubling",
		Cost:      "1M + 3S + 1*roota6",
	},
	{
		ID:        "g12o/shortw/xz/doubling/dbl-2003-s-2",
		Operation: "doubling",
		Cost:      "1M + 4S + 1*a6",
	},
	{
		ID:        "g12o/shortw/xz/doubling/dbl-2003-s",
		Operation: "doubling",
		Cost:      "1M + 1S + 2^4 + 1*a6",
	},
	{
		ID:        "g12o/shortw/xz/diffadd/mdadd-2003-s",
		Operation: "diffadd",
		Cost:      "4M + 1S",
	},
	{
		ID:        "g12o/shortw/xz/diffadd/mdadd-2003-s-2",
		Operation: "diffadd",
		Cost:      "4M + 3S + 1*a6",
	},
	{
		ID:        "g12o/shortw/xz/diffadd/dadd-2003-s-2",
		Operation: "diffadd",
		Cost:      "5M + 3S + 1*a6",
	},
	{
		ID:        "g12o/shortw/xz/diffadd/dadd-2003-s",
		Operation: "diffadd",
		Cost:      "7M + 5S + 1*a6",
	},
	{
		ID:        "g12o/shortw/xz/ladder/mladd-2003-s-2",
		Operation: "ladder",
		Cost:      "5M + 4S + 1*sqrta6",
	},
	{
		ID:        "g12o/shortw/xz/ladder/mladd-2003-s",
		Operation: "ladder",
		Cost:      "5M + 5S + 1*a6",
	},
	{
		ID:        "g12o/shortw/xz/ladder/mladd-2003-s-3",
		Operation: "ladder",
		Cost:      "5M + 5S + 2*sqrta6",
	},
	{
		ID:        "g12o/shortw/xz/ladder/ladd-2003-s-3",
		Operation: "ladder",
		Cost:      "6M + 5S + 2*sqrta6",
	},
	{
		ID:        "g12o/shortw/xz/ladder/ladd-2003-s-4",
		Operation: "ladder",
		Cost:      "6M + 5S + 1*roota6 + 1*sqrta6",
	},
	{
		ID:        "g12o/shortw/xz/ladder/ladd-2003-s-2",
		Operation: "ladder",
		Cost:      "6M + 7S + 2*a6",
	},
	{
		ID:        "g12o/shortw/xz/ladder/ladd-2003-s",
		Operation: "ladder",
		Cost:      "8M + 6S + 2^4 + 2*a6",
	},
	{
		ID:        "g12o/shortw/xz/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 1M",
	},
	{
		ID:        "g1p/2dik/standard/addition/mmadd-20080313-bl",
		Operation: "addition",
		Cost:      "4M + 4S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/mmadd-20080313-bl",
		Operation: "readdition",
		Cost:      "4M + 4S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/mmadd-20080308-bl",
		Operation: "addition",
		Cost:      "4M + 4S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/mmadd-20080308-bl",
		Operation: "readdition",
		Cost:      "4M + 4S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/mmadd-2006-dik",
		Operation: "addition",
		Cost:      "6M + 3S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/mmadd-2006-dik",
		Operation: "readdition",
		Cost:      "6M + 3S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/madd-2007-bl",
		Operation: "addition",
		Cost:      "8M + 4S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/madd-2007-bl",
		Operation: "readdition",
		Cost:      "8M + 4S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/madd-2006-dik",
		Operation: "addition",
		Cost:      "9M + 3S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/madd-2006-dik",
		Operation: "readdition",
		Cost:      "9M + 3S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/add-2006-dik-3",
		Operation: "addition",
		Cost:      "12M + 5S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/add-2006-dik-3",
		Operation: "readdition",
		Cost:      "12M + 5S + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/add-2006-dik-2",
		Operation: "addition",
		Cost:      "21M + 15S + 2^4 + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/add-2006-dik-2",
		Operation: "readdition",
		Cost:      "21M + 11S + 1^4 + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/add-2006-dik",
		Operation: "addition",
		Cost:      "7I + 12M + 9S + 1^4 + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/addition/add-2006-dik",
		Operation: "readdition",
		Cost:      "4I + 9M + 8S + 1^4 + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 5S + 1*a2 + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*a2 + 1*a",
	},
	{
		ID:        "g1p/2dik/standard/doubling/dbl-2006-dik-2",
		Operation: "doubling",
		Cost:      "3M + 4S + 1*a + 1*a4",
	},
	{
		ID:        "g1p/2dik/standard/doubling/dbl-2006-dik",
		Operation: "doubling",
		Cost:      "3M + 8S + 1*a16 + 2*a",
	},
	{
		ID:        "g1p/2dik/standard/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M + 1S",
	},
	{
		ID:        "g1p/3dik/standard/addition/mmadd-2007-bblp",
		Operation: "addition",
		Cost:      "4M + 2S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/mmadd-2007-bblp",
		Operation: "readdition",
		Cost:      "4M + 2S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/mmadd-2006-dik",
		Operation: "addition",
		Cost:      "4M + 3S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/mmadd-2006-dik",
		Operation: "readdition",
		Cost:      "4M + 3S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/madd-2007-bblp",
		Operation: "addition",
		Cost:      "7M + 4S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/madd-2007-bblp",
		Operation: "readdition",
		Cost:      "7M + 4S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/madd-2006-dik",
		Operation: "addition",
		Cost:      "8M + 3S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/madd-2006-dik",
		Operation: "readdition",
		Cost:      "8M + 3S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/add-2007-bblp",
		Operation: "addition",
		Cost:      "11M + 6S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/add-2007-bblp",
		Operation: "readdition",
		Cost:      "10M + 6S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/add-2006-dik-3",
		Operation: "addition",
		Cost:      "13M + 4S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/add-2006-dik-3",
		Operation: "readdition",
		Cost:      "12M + 4S + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/addition/add-2006-dik-2",
		Operation: "addition",
		Cost:      "15M + 11S + 5^3 + 1*a",
	},
	{
		ID:        "g1p/3dik/standard/addition/add-2006-dik-2",
		Operation: "readdition",
		Cost:      "15M + 7S + 3^3 + 1*a",
	},
	{
		ID:        "g1p/3dik/standard/addition/add-2006-dik",
		Operation: "addition",
		Cost:      "7I + 11M + 9S + 5^3 + 1*a",
	},
	{
		ID:        "g1p/3dik/standard/addition/add-2006-dik",
		Operation: "readdition",
		Cost:      "5I + 9M + 8S + 4^3 + 1*a",
	},
	{
		ID:        "g1p/3dik/standard/doubling/mdbl-2007-bblp",
		Operation: "doubling",
		Cost:      "1M + 5S + 1*a2 + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/doubling/dbl-2007-bblp",
		Operation: "doubling",
		Cost:      "2M + 7S + 1*a2 + 1*a3",
	},
	{
		ID:        "g1p/3dik/standard/doubling/dbl-2006-dik-2",
		Operation: "doubling",
		Cost:      "4M + 5S + 1*a3 + 1*a",
	},
	{
		ID:        "g1p/3dik/standard/doubling/dbl-2006-dik",
		Operation: "doubling",
		Cost:      "4M + 5S + 1^4 + 1*a3 + 2*a",
	},
	{
		ID:        "g1p/3dik/standard/tripling/tpl-2006-dik-2",
		Operation: "tripling",
		Cost:      "6M + 6S + 2*a",
	},
	{
		ID:        "g1p/3dik/standard/tripling/tpl-2006-dik",
		Operation: "tripling",
		Cost:      "6M + 7S + 1*a + 1*b + 1*c",
	},
	{
		ID:        "g1p/3dik/standard/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 3M + 1S",
	},
	{
		ID:        "g1p/edwards/inverted/addition/mmadd-2007-bl",
		Operation: "addition",
		Cost:      "7M + 2*c",
	},
	{
		ID:        "g1p/edwards/inverted/addition/mmadd-2007-bl",
		Operation: "readdition",
		Cost:      "7M + 2*c",
	},
	{
		ID:        "g1p/edwards/inverted/addition/xmadd-2007-bl",
		Operation: "addition",
		Cost:      "8M + 1S + 2*c + 1*d",
	},
	{
		ID:        "g1p/edwards/inverted/addition/xmadd-2007-bl",
		Operation: "readdition",
		Cost:      "8M + 1S + 2*c + 1*d",
	},
	{
		ID:        "g1p/edwards/inverted/addition/madd-2007-bl",
		Operation: "addition",
		Cost:      "8M + 1S + 2*c + 1*d",
	},
	{
		ID:        "g1p/edwards/inverted/addition/madd-2007-bl",
		Operation: "readdition",
		Cost:      "8M + 1S + 2*c + 1*d",
	},
	{
		ID:        "g1p/edwards/inverted/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "9M + 1*c",
	},
	{
		ID:        "g1p/edwards/inverted/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "9M + 1*c",
	},
	{
		ID:        "g1p/edwards/inverted/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "9M + 1S + 2*c + 1*d",
	},
	{
		ID:        "g1p/edwards/inverted/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "9M + 1S + 2*c + 1*d",
	},
	{
		ID:        "g1p/edwards/inverted/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "11M + 1*c",
	},
	{
		ID:        "g1p/edwards/inverted/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "11M + 1*c",
	},
	{
		ID:        "g1p/edwards/inverted/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "3M + 3S + 1*c",
	},
	{
		ID:        "g1p/edwards/inverted/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "3M + 4S + 1*ccd2 + 1*c",
	},
	{
		ID:        "g1p/edwards/inverted/tripling/tpl-2007-bl",
		Operation: "tripling",
		Cost:      "9M + 4S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/inverted/tripling/tpl-2007-bl-2",
		Operation: "tripling",
		Cost:      "7M + 7S + 1*ccd",
	},
	{
		ID:        "g1p/edwards/inverted/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g1p/edwards/projective/addition/mmadd-2007-bl",
		Operation: "addition",
		Cost:      "6M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/mmadd-2007-bl",
		Operation: "readdition",
		Cost:      "6M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "9M + 1*k",
	},
	{
		ID:        "g1p/edwards/projective/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "9M + 1*k",
	},
	{
		ID:        "g1p/edwards/projective/addition/xmadd-2007-hcd",
		Operation: "addition",
		Cost:      "9M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/xmadd-2007-hcd",
		Operation: "readdition",
		Cost:      "9M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/madd-2007-bl-2",
		Operation: "addition",
		Cost:      "9M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/madd-2007-bl-2",
		Operation: "readdition",
		Cost:      "9M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/madd-2007-bl",
		Operation: "addition",
		Cost:      "9M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/madd-2007-bl",
		Operation: "readdition",
		Cost:      "9M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/madd-2007-bl-3",
		Operation: "addition",
		Cost:      "6M + 5S + 1*c2 + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/madd-2007-bl-3",
		Operation: "readdition",
		Cost:      "6M + 5S + 1*c2 + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-2007-bl-2",
		Operation: "addition",
		Cost:      "10M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-2007-bl-2",
		Operation: "readdition",
		Cost:      "10M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "10M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "10M + 1S + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-2007-bl-4",
		Operation: "addition",
		Cost:      "10M + 1S + 3*i + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-2007-bl-4",
		Operation: "readdition",
		Cost:      "10M + 1S + 2*i + 1*c + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "11M + 1*k",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "11M + 1*k",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-2007-bl-3",
		Operation: "addition",
		Cost:      "7M + 5S + 1*c2 + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-2007-bl-3",
		Operation: "readdition",
		Cost:      "7M + 5S + 1*c2 + 1*d",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-20090311-hwcd",
		Operation: "addition",
		Cost:      "10M + 3S + 1*k",
	},
	{
		ID:        "g1p/edwards/projective/addition/add-20090311-hwcd",
		Operation: "readdition",
		Cost:      "9M + 2S + 1*k",
	},
	{
		ID:        "g1p/edwards/projective/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "3M + 3S + 2*c",
	},
	{
		ID:        "g1p/edwards/projective/doubling/dbl-2007-bl-2",
		Operation: "doubling",
		Cost:      "3M + 4S + 3*c",
	},
	{
		ID:        "g1p/edwards/projective/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "3M + 4S + 3*c",
	},
	{
		ID:        "g1p/edwards/projective/doubling/dbl-2007-bl-3",
		Operation: "doubling",
		Cost:      "3M + 4S + 3*c",
	},
	{
		ID:        "g1p/edwards/projective/tripling/tpl-2007-bblp",
		Operation: "tripling",
		Cost:      "9M + 4S + 1*c2",
	},
	{
		ID:        "g1p/edwards/projective/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "9M + 4S + 1*c",
	},
	{
		ID:        "g1p/edwards/projective/tripling/tpl-2007-bblp-2",
		Operation: "tripling",
		Cost:      "7M + 7S",
	},
	{
		ID:        "g1p/edwards/projective/tripling/tpl-2007-bblp-3",
		Operation: "tripling",
		Cost:      "7M + 7S + 1*cc4",
	},
	{
		ID:        "g1p/edwards/projective/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g1p/edwards/yz/doubling/mdbl-2006-g-2",
		Operation: "doubling",
		Cost:      "2S + 1*r2 + 1*s",
	},
	{
		ID:        "g1p/edwards/yz/doubling/mdbl-2006-g-3",
		Operation: "doubling",
		Cost:      "3S + 1*s",
	},
	{
		ID:        "g1p/edwards/yz/doubling/dbl-2006-g-2",
		Operation: "doubling",
		Cost:      "4S + 1*r + 1*s",
	},
	{
		ID:        "g1p/edwards/yz/doubling/dbl-2006-g",
		Operation: "doubling",
		Cost:      "6S + 2*r + 1*s",
	},
	{
		ID:        "g1p/edwards/yz/diffadd/mdadd-2006-g-2",
		Operation: "diffadd",
		Cost:      "3M + 4S + 3*r + 1*s",
	},
	{
		ID:        "g1p/edwards/yz/diffadd/dadd-2006-g-2",
		Operation: "diffadd",
		Cost:      "4M + 4S + 3*r + 1*s",
	},
	{
		ID:        "g1p/edwards/yz/diffadd/dadd-2006-g",
		Operation: "diffadd",
		Cost:      "4M + 8S + 5*r + 1*s",
	},
	{
		ID:        "g1p/edwards/yz/ladder/mladd-2006-g-2",
		Operation: "ladder",
		Cost:      "3M + 6S + 3*r + 2*s",
	},
	{
		ID:        "g1p/edwards/yz/ladder/ladd-2006-g-2",
		Operation: "ladder",
		Cost:      "4M + 6S + 3*r + 2*s",
	},
	{
		ID:        "g1p/edwards/yz/ladder/ladd-2006-g",
		Operation: "ladder",
		Cost:      "4M + 14S + 7*r + 2*s",
	},
	{
		ID:        "g1p/edwards/yz/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 1M",
	},
	{
		ID:        "g1p/edwards/yzsquared/doubling/mdbl-2006-g",
		Operation: "doubling",
		Cost:      "3S + 1*r + 1*s",
	},
	{
		ID:        "g1p/edwards/yzsquared/doubling/dbl-2006-g",
		Operation: "doubling",
		Cost:      "4S + 1*r + 1*s",
	},
	{
		ID:        "g1p/edwards/yzsquared/diffadd/mdadd-2006-g",
		Operation: "diffadd",
		Cost:      "3M + 2S + 1*s",
	},
	{
		ID:        "g1p/edwards/yzsquared/diffadd/dadd-2006-g",
		Operation: "diffadd",
		Cost:      "4M + 2S + 1*s",
	},
	{
		ID:        "g1p/edwards/yzsquared/ladder/mladd-2006-g-2",
		Operation: "ladder",
		Cost:      "3M + 6S + 1*r + 2*s",
	},
	{
		ID:        "g1p/edwards/yzsquared/ladder/mladd-2006-g",
		Operation: "ladder",
		Cost:      "3M + 6S + 1*r + 2*s",
	},
	{
		ID:        "g1p/edwards/yzsquared/ladder/ladd-2006-g-2",
		Operation: "ladder",
		Cost:      "4M + 6S + 1*r + 2*s",
	},
	{
		ID:        "g1p/edwards/yzsquared/ladder/ladd-2006-g",
		Operation: "ladder",
		Cost:      "4M + 6S + 1*r + 2*s",
	},
	{
		ID:        "g1p/edwards/yzsquared/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 1M",
	},
	{
		ID:        "g1p/hessian/extended/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "5M + 6S",
	},
	{
		ID:        "g1p/hessian/extended/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "5M + 6S",
	},
	{
		ID:        "g1p/hessian/extended/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "6M + 6S",
	},
	{
		ID:        "g1p/hessian/extended/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "6M + 6S",
	},
	{
		ID:        "g1p/hessian/extended/doubling/dbl-20080225-hwcd",
		Operation: "doubling",
		Cost:      "3M + 6S",
	},
	{
		ID:        "g1p/hessian/extended/doubling/mdbl-20080225-hwcd",
		Operation: "doubling",
		Cost:      "3M + 6S",
	},
	{
		ID:        "g1p/hessian/extended/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 3M + 2S",
	},
	{
		ID:        "g1p/hessian/standard/addition/mmadd",
		Operation: "addition",
		Cost:      "8M",
	},
	{
		ID:        "g1p/hessian/standard/addition/mmadd",
		Operation: "readdition",
		Cost:      "7M",
	},
	{
		ID:        "g1p/hessian/standard/addition/madd",
		Operation: "addition",
		Cost:      "10M",
	},
	{
		ID:        "g1p/hessian/standard/addition/madd",
		Operation: "readdition",
		Cost:      "10M",
	},
	{
		ID:        "g1p/hessian/standard/addition/add-2001-jq",
		Operation: "addition",
		Cost:      "12M",
	},
	{
		ID:        "g1p/hessian/standard/addition/add-2001-jq",
		Operation: "readdition",
		Cost:      "12M",
	},
	{
		ID:        "g1p/hessian/standard/addition/add-2009-bkl",
		Operation: "addition",
		Cost:      "12M",
	},
	{
		ID:        "g1p/hessian/standard/addition/add-2009-bkl",
		Operation: "readdition",
		Cost:      "12M",
	},
	{
		ID:        "g1p/hessian/standard/addition/add2",
		Operation: "addition",
		Cost:      "12M",
	},
	{
		ID:        "g1p/hessian/standard/addition/add2",
		Operation: "readdition",
		Cost:      "12M",
	},
	{
		ID:        "g1p/hessian/standard/addition/xmadd-2007-hcd",
		Operation: "addition",
		Cost:      "5M + 9S",
	},
	{
		ID:        "g1p/hessian/standard/addition/xmadd-2007-hcd",
		Operation: "readdition",
		Cost:      "5M + 6S",
	},
	{
		ID:        "g1p/hessian/standard/addition/add-2008-hwcd",
		Operation: "addition",
		Cost:      "6M + 12S",
	},
	{
		ID:        "g1p/hessian/standard/addition/add-2008-hwcd",
		Operation: "readdition",
		Cost:      "6M + 6S",
	},
	{
		ID:        "g1p/hessian/standard/addition/add",
		Operation: "addition",
		Cost:      "12M + 6S",
	},
	{
		ID:        "g1p/hessian/standard/addition/add",
		Operation: "readdition",
		Cost:      "9M + 3S",
	},
	{
		ID:        "g1p/hessian/standard/doubling/mdbl-2007-hcd",
		Operation: "doubling",
		Cost:      "3M + 3S",
	},
	{
		ID:        "g1p/hessian/standard/doubling/dbl-2007-hcd-2",
		Operation: "doubling",
		Cost:      "7M + 1S",
	},
	{
		ID:        "g1p/hessian/standard/doubling/dbl-2007-hcd",
		Operation: "doubling",
		Cost:      "7M + 1S",
	},
	{
		ID:        "g1p/hessian/standard/doubling/dbl-2007-hcd-3",
		Operation: "doubling",
		Cost:      "3M + 6S",
	},
	{
		ID:        "g1p/hessian/standard/doubling/dbl-2007-hcd-4",
		Operation: "doubling",
		Cost:      "3M + 6S",
	},
	{
		ID:        "g1p/hessian/standard/doubling/dbl2",
		Operation: "doubling",
		Cost:      "6M + 3S",
	},
	{
		ID:        "g1p/hessian/standard/doubling/dbl-2001-jq",
		Operation: "doubling",
		Cost:      "12M",
	},
	{
		ID:        "g1p/hessian/standard/doubling/dbl",
		Operation: "doubling",
		Cost:      "3M + 6^3",
	},
	{
		ID:        "g1p/hessian/standard/tripling/tpl-2007-hcd-3",
		Operation: "tripling",
		Cost:      "8M + 6S + 1*b",
	},
	{
		ID:        "g1p/hessian/standard/tripling/tpl-2007-hcd-2",
		Operation: "tripling",
		Cost:      "11M + 4S + 2*a",
	},
	{
		ID:        "g1p/hessian/standard/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "10M + 1S + 29^3 + 2*d",
	},
	{
		ID:        "g1p/hessian/standard/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g1p/jintersect/extended/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "10M + 1S + 2*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "10M + 1S + 2*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/mmadd-2001-ls",
		Operation: "addition",
		Cost:      "10M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/mmadd-2001-ls",
		Operation: "readdition",
		Cost:      "10M + 1S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "11M + 1S + 2*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "11M + 1S + 2*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/madd-2001-ls",
		Operation: "addition",
		Cost:      "13M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/madd-2001-ls",
		Operation: "readdition",
		Cost:      "12M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/smadd-2001-ls",
		Operation: "addition",
		Cost:      "13M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/smadd-2001-ls",
		Operation: "readdition",
		Cost:      "12M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/add-2001-ls",
		Operation: "addition",
		Cost:      "15M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/add-2001-ls",
		Operation: "readdition",
		Cost:      "13M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/add-1986-cc-2",
		Operation: "addition",
		Cost:      "16M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/add-1986-cc-2",
		Operation: "readdition",
		Cost:      "14M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/add-1986-cc",
		Operation: "addition",
		Cost:      "22M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/addition/add-1986-cc",
		Operation: "readdition",
		Cost:      "20M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/doubling/dbl-20080225-hwcd",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*a",
	},
	{
		ID:        "g1p/jintersect/extended/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "4M + 4S",
	},
	{
		ID:        "g1p/jintersect/extended/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "5M + 4S",
	},
	{
		ID:        "g1p/jintersect/extended/doubling/dbl-2001-ls",
		Operation: "doubling",
		Cost:      "6M + 3S",
	},
	{
		ID:        "g1p/jintersect/extended/doubling/dbl-1986-cc-2",
		Operation: "doubling",
		Cost:      "7M + 3S",
	},
	{
		ID:        "g1p/jintersect/extended/doubling/dbl-1986-cc",
		Operation: "doubling",
		Cost:      "14M + 9S",
	},
	{
		ID:        "g1p/jintersect/extended/tripling/tpl-2007-hcd-4",
		Operation: "tripling",
		Cost:      "6M + 10S + 1*bb2 + 1*b2 + 2*a + 1*b3",
	},
	{
		ID:        "g1p/jintersect/extended/tripling/tpl-2007-hcd-3",
		Operation: "tripling",
		Cost:      "6M + 10S + 1*bb2 + 1*b2 + 2*a + 1*b3",
	},
	{
		ID:        "g1p/jintersect/extended/tripling/tpl-2007-hcd-2",
		Operation: "tripling",
		Cost:      "9M + 7S + 3*b",
	},
	{
		ID:        "g1p/jintersect/extended/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "9M + 7S + 5*b",
	},
	{
		ID:        "g1p/jintersect/extended/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 4M",
	},
	{
		ID:        "g1p/jintersect/standard/addition/mmadd-2001-ls",
		Operation: "addition",
		Cost:      "8M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/mmadd-2001-ls",
		Operation: "readdition",
		Cost:      "8M + 1S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "11M + 1S + 2*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "10M + 1S + 2*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/madd-2001-ls",
		Operation: "addition",
		Cost:      "11M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/madd-2001-ls",
		Operation: "readdition",
		Cost:      "10M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/smadd-2001-ls",
		Operation: "addition",
		Cost:      "11M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/smadd-2001-ls",
		Operation: "readdition",
		Cost:      "10M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "13M + 1S + 2*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "11M + 1S + 2*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/add-2001-ls",
		Operation: "addition",
		Cost:      "13M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/add-2001-ls",
		Operation: "readdition",
		Cost:      "11M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/add-1986-cc-2",
		Operation: "addition",
		Cost:      "14M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/add-1986-cc-2",
		Operation: "readdition",
		Cost:      "12M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/add-1986-cc",
		Operation: "addition",
		Cost:      "20M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/addition/add-1986-cc",
		Operation: "readdition",
		Cost:      "18M + 2S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/doubling/mdbl-20090427-b",
		Operation: "doubling",
		Cost:      "6S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/doubling/mdbl-20080225-hwcd",
		Operation: "doubling",
		Cost:      "1M + 5S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "2M + 4S",
	},
	{
		ID:        "g1p/jintersect/standard/doubling/dbl-20080225-hwcd",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*a",
	},
	{
		ID:        "g1p/jintersect/standard/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "3M + 4S",
	},
	{
		ID:        "g1p/jintersect/standard/doubling/dbl-2001-ls",
		Operation: "doubling",
		Cost:      "4M + 3S",
	},
	{
		ID:        "g1p/jintersect/standard/doubling/dbl-1986-cc-2",
		Operation: "doubling",
		Cost:      "5M + 3S",
	},
	{
		ID:        "g1p/jintersect/standard/doubling/dbl-1986-cc",
		Operation: "doubling",
		Cost:      "12M + 9S",
	},
	{
		ID:        "g1p/jintersect/standard/tripling/tpl-2007-hcd-4",
		Operation: "tripling",
		Cost:      "4M + 10S + 1*bb2 + 1*b2 + 2*a + 1*b3",
	},
	{
		ID:        "g1p/jintersect/standard/tripling/tpl-2007-hcd-3",
		Operation: "tripling",
		Cost:      "4M + 10S + 1*bb2 + 1*b2 + 2*a + 1*b3",
	},
	{
		ID:        "g1p/jintersect/standard/tripling/tpl-2007-hcd-2",
		Operation: "tripling",
		Cost:      "7M + 7S + 3*b",
	},
	{
		ID:        "g1p/jintersect/standard/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "7M + 7S + 5*b",
	},
	{
		ID:        "g1p/jintersect/standard/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 3M",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "6M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "6M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "7M + 4S + 1*k",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "7M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/doubling/mdbl-20090311-hwcd",
		Operation: "doubling",
		Cost:      "6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/doubling/mdbl-2007-hcd",
		Operation: "doubling",
		Cost:      "1M + 5S",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/doubling/dbl-20090311-hwcd",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/doubling/dbl-2007-hcd",
		Operation: "doubling",
		Cost:      "3M + 4S",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/doubling/dbl-2007-fw-4",
		Operation: "doubling",
		Cost:      "8S + 1*a + 2*c",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/doubling/dbl-2007-fw-2",
		Operation: "doubling",
		Cost:      "1M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/doubling/dbl-2007-fw",
		Operation: "doubling",
		Cost:      "3M + 8S + 1*a2 + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/doubling/dbl-2007-fw-3",
		Operation: "doubling",
		Cost:      "2M + 7S + 1^4 + 1*a2 + 2*c",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "8M + 6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/tripling/tpl-2007-hcd-2",
		Operation: "tripling",
		Cost:      "4M + 11S + 1*a + 1*b",
	},
	{
		ID:        "g1p/jquartic/2xxyzz/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M + 2S",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/mmadd-2007-bl",
		Operation: "addition",
		Cost:      "5M + 4S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/mmadd-2007-bl",
		Operation: "readdition",
		Cost:      "5M + 4S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "6M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "6M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/madd-2007-bl",
		Operation: "addition",
		Cost:      "7M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/madd-2007-bl",
		Operation: "readdition",
		Cost:      "7M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "7M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "7M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/add-2007-d",
		Operation: "addition",
		Cost:      "9M + 2S + 2*half + 1*b",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/addition/add-2007-d",
		Operation: "readdition",
		Cost:      "9M + 2S + 1*half + 1*b",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/doubling/mdbl-20090311-hwcd",
		Operation: "doubling",
		Cost:      "1M + 6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/doubling/mdbl-2007-hcd",
		Operation: "doubling",
		Cost:      "1M + 6S",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/doubling/dbl-20090311-hwcd",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/doubling/dbl-2007-hcd",
		Operation: "doubling",
		Cost:      "3M + 4S",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/doubling/dbl-2007-fw-4",
		Operation: "doubling",
		Cost:      "8S + 1*a + 2*c",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/doubling/dbl-2007-fw-2",
		Operation: "doubling",
		Cost:      "1M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/doubling/dbl-2007-fw",
		Operation: "doubling",
		Cost:      "3M + 9S + 2*a2",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/doubling/dbl-2007-fw-3",
		Operation: "doubling",
		Cost:      "2M + 8S + 1^4 + 1*a + 2*c",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/tripling/tpl-2007-hcd-2",
		Operation: "tripling",
		Cost:      "4M + 11S + 1*a + 1*b",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "8M + 7S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xxyzzr/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M + 1S",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/mmadd-2002-bj-2",
		Operation: "addition",
		Cost:      "5M + 2S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/mmadd-2002-bj-2",
		Operation: "readdition",
		Cost:      "5M + 2S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/madd-2002-bj",
		Operation: "addition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/madd-2002-bj",
		Operation: "readdition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2002-bj-2",
		Operation: "addition",
		Cost:      "10M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2002-bj-2",
		Operation: "readdition",
		Cost:      "9M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "8M + 6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2002-bj-3",
		Operation: "addition",
		Cost:      "10M + 4S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2002-bj-3",
		Operation: "readdition",
		Cost:      "9M + 2S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2007-d",
		Operation: "addition",
		Cost:      "10M + 4S + 1*b",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2007-d",
		Operation: "readdition",
		Cost:      "9M + 2S + 1*b",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2002-bj",
		Operation: "addition",
		Cost:      "19M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/add-2002-bj",
		Operation: "readdition",
		Cost:      "18M + 6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/mmadd-2002-bj",
		Operation: "addition",
		Cost:      "2I + 11M + 5S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/addition/mmadd-2002-bj",
		Operation: "readdition",
		Cost:      "2I + 11M + 4S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/mdbl-2007-fw",
		Operation: "doubling",
		Cost:      "1M + 4S + 1*a2",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/dbl-2007-fw-4",
		Operation: "doubling",
		Cost:      "1M + 7S + 1*a2 + 1*c2 + 1*c",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/dbl-2007-hcd",
		Operation: "doubling",
		Cost:      "2M + 6S + 1*a + 1*b",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/dbl-2007-fw-2",
		Operation: "doubling",
		Cost:      "2M + 6S + 1*a2",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/dbl-2007-fw",
		Operation: "doubling",
		Cost:      "3M + 6S + 2*a2",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 9S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/dbl-2007-fw-3",
		Operation: "doubling",
		Cost:      "2M + 5S + 1^4 + 1*a2 + 2*c",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/dbl-2002-bj-2",
		Operation: "doubling",
		Cost:      "19M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/doubling/dbl-2002-bj",
		Operation: "doubling",
		Cost:      "19M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/2xyz/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M + 1S",
	},
	{
		ID:        "g1p/jquartic/xxyzz/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "6M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/xxyzz/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "6M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/xxyzz/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "7M + 4S + 1*k",
	},
	{
		ID:        "g1p/jquartic/xxyzz/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "7M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/xxyzz/doubling/mdbl-20090311-hwcd",
		Operation: "doubling",
		Cost:      "6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzz/doubling/mdbl-2007-hcd",
		Operation: "doubling",
		Cost:      "1M + 5S",
	},
	{
		ID:        "g1p/jquartic/xxyzz/doubling/dbl-20090311-hwcd",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzz/doubling/dbl-2007-hcd",
		Operation: "doubling",
		Cost:      "3M + 4S",
	},
	{
		ID:        "g1p/jquartic/xxyzz/doubling/dbl-2007-fw-2",
		Operation: "doubling",
		Cost:      "1M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzz/doubling/dbl-2007-fw",
		Operation: "doubling",
		Cost:      "3M + 8S + 1*a2 + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzz/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "8M + 6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzz/tripling/tpl-2007-hcd-2",
		Operation: "tripling",
		Cost:      "4M + 11S + 1*a + 1*b",
	},
	{
		ID:        "g1p/jquartic/xxyzz/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M + 2S",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/mmadd-2007-bl",
		Operation: "addition",
		Cost:      "5M + 4S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/mmadd-2007-bl",
		Operation: "readdition",
		Cost:      "5M + 4S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/madd-20080225-hwcd",
		Operation: "addition",
		Cost:      "6M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/madd-20080225-hwcd",
		Operation: "readdition",
		Cost:      "6M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/madd-2007-bl",
		Operation: "addition",
		Cost:      "7M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/madd-2007-bl",
		Operation: "readdition",
		Cost:      "7M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/add-20080225-hwcd",
		Operation: "addition",
		Cost:      "7M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/add-20080225-hwcd",
		Operation: "readdition",
		Cost:      "7M + 3S + 1*k",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/add-2007-d",
		Operation: "addition",
		Cost:      "9M + 2S + 2*half + 1*b",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/addition/add-2007-d",
		Operation: "readdition",
		Cost:      "9M + 2S + 1*half + 1*b",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/doubling/mdbl-20090311-hwcd",
		Operation: "doubling",
		Cost:      "1M + 6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/doubling/mdbl-2007-hcd",
		Operation: "doubling",
		Cost:      "1M + 6S",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/doubling/dbl-20090311-hwcd",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/doubling/dbl-2007-hcd",
		Operation: "doubling",
		Cost:      "3M + 4S",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/doubling/dbl-2007-fw-2",
		Operation: "doubling",
		Cost:      "1M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/doubling/dbl-2007-fw",
		Operation: "doubling",
		Cost:      "3M + 9S + 2*a2",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/tripling/tpl-2007-hcd-2",
		Operation: "tripling",
		Cost:      "4M + 11S + 1*a + 1*b",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/tripling/tpl-2007-hcd",
		Operation: "tripling",
		Cost:      "8M + 7S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xxyzzr/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M + 1S",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/mmadd-2002-bj-2",
		Operation: "addition",
		Cost:      "5M + 2S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/mmadd-2002-bj-2",
		Operation: "readdition",
		Cost:      "5M + 2S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/madd-2002-bj",
		Operation: "addition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/madd-2002-bj",
		Operation: "readdition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2002-bj-2",
		Operation: "addition",
		Cost:      "10M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2002-bj-2",
		Operation: "readdition",
		Cost:      "9M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "8M + 6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "8M + 3S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2002-bj-3",
		Operation: "addition",
		Cost:      "10M + 4S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2002-bj-3",
		Operation: "readdition",
		Cost:      "9M + 2S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2007-d",
		Operation: "addition",
		Cost:      "10M + 4S + 1*b",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2007-d",
		Operation: "readdition",
		Cost:      "9M + 2S + 1*b",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2002-bj",
		Operation: "addition",
		Cost:      "19M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/add-2002-bj",
		Operation: "readdition",
		Cost:      "18M + 6S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/mmadd-2002-bj",
		Operation: "addition",
		Cost:      "2I + 11M + 5S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/addition/mmadd-2002-bj",
		Operation: "readdition",
		Cost:      "2I + 11M + 4S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/doubling/mdbl-2007-fw",
		Operation: "doubling",
		Cost:      "1M + 4S + 1*a2",
	},
	{
		ID:        "g1p/jquartic/xyz/doubling/dbl-2007-hcd",
		Operation: "doubling",
		Cost:      "2M + 6S + 1*a + 1*b",
	},
	{
		ID:        "g1p/jquartic/xyz/doubling/dbl-2007-fw-2",
		Operation: "doubling",
		Cost:      "2M + 6S + 1*a2",
	},
	{
		ID:        "g1p/jquartic/xyz/doubling/dbl-2007-fw",
		Operation: "doubling",
		Cost:      "3M + 6S + 2*a2",
	},
	{
		ID:        "g1p/jquartic/xyz/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 9S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/doubling/dbl-2002-bj-2",
		Operation: "doubling",
		Cost:      "19M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/doubling/dbl-2002-bj",
		Operation: "doubling",
		Cost:      "19M + 8S + 1*a",
	},
	{
		ID:        "g1p/jquartic/xyz/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M + 1S",
	},
	{
		ID:        "g1p/montgom/xz/doubling/mdbl-1987-m",
		Operation: "doubling",
		Cost:      "1M + 2S + 1*a",
	},
	{
		ID:        "g1p/montgom/xz/doubling/dbl-1987-m-3",
		Operation: "doubling",
		Cost:      "2M + 2S + 1*a24",
	},
	{
		ID:        "g1p/montgom/xz/doubling/dbl-1987-m-2",
		Operation: "doubling",
		Cost:      "4M + 3S + 1*a24",
	},
	{
		ID:        "g1p/montgom/xz/doubling/dbl-1987-m",
		Operation: "doubling",
		Cost:      "3M + 5S + 1*a",
	},
	{
		ID:        "g1p/montgom/xz/diffadd/mdadd-1987-m",
		Operation: "diffadd",
		Cost:      "3M + 2S",
	},
	{
		ID:        "g1p/montgom/xz/diffadd/dadd-1987-m-3",
		Operation: "diffadd",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/montgom/xz/diffadd/dadd-1987-m",
		Operation: "diffadd",
		Cost:      "6M + 2S",
	},
	{
		ID:        "g1p/montgom/xz/diffadd/dadd-1987-m-2",
		Operation: "diffadd",
		Cost:      "6M + 2S",
	},
	{
		ID:        "g1p/montgom/xz/ladder/mladd-1987-m",
		Operation: "ladder",
		Cost:      "5M + 4S + 1*a24",
	},
	{
		ID:        "g1p/montgom/xz/ladder/ladd-1987-m-3",
		Operation: "ladder",
		Cost:      "6M + 4S + 1*a24",
	},
	{
		ID:        "g1p/montgom/xz/ladder/ladd-1987-m-2",
		Operation: "ladder",
		Cost:      "10M + 5S + 1*a24",
	},
	{
		ID:        "g1p/montgom/xz/ladder/ladd-1987-m",
		Operation: "ladder",
		Cost:      "9M + 7S + 1*a",
	},
	{
		ID:        "g1p/montgom/xz/scaling/scale",
		Operation: "scaling",
		Cost:      "1I + 1M",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/mmadd-2007-bl",
		Operation: "addition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/mmadd-2007-bl",
		Operation: "readdition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/zadd-2007-m",
		Operation: "addition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/zadd-2007-m",
		Operation: "readdition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/madd-2007-bl",
		Operation: "addition",
		Cost:      "7M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/madd-2007-bl",
		Operation: "readdition",
		Cost:      "7M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/madd-2004-hmv",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/madd-2004-hmv",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/madd-2008-g",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/madd-2008-g",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/madd",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/madd",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "11M + 5S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "10M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1998-cmo-2",
		Operation: "addition",
		Cost:      "12M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1998-cmo-2",
		Operation: "readdition",
		Cost:      "11M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-2001-b",
		Operation: "addition",
		Cost:      "12M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-2001-b",
		Operation: "readdition",
		Cost:      "11M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1998-hnm",
		Operation: "addition",
		Cost:      "12M + 4S + 1*half",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1998-hnm",
		Operation: "readdition",
		Cost:      "11M + 3S + 1*half",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1986-cc-2",
		Operation: "addition",
		Cost:      "8M + 6S + 2^3",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1986-cc-2",
		Operation: "readdition",
		Cost:      "8M + 5S + 1^3",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1986-cc",
		Operation: "addition",
		Cost:      "10M + 5S + 3^3",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1986-cc",
		Operation: "readdition",
		Cost:      "10M + 4S + 2^3",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1998-cmo",
		Operation: "addition",
		Cost:      "10M + 5S + 4^3",
	},
	{
		ID:        "g1p/shortw/jacobian-0/addition/add-1998-cmo",
		Operation: "readdition",
		Cost:      "10M + 4S + 3^3",
	},
	{
		ID:        "g1p/shortw/jacobian-0/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 5S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/doubling/dbl-2009-l",
		Operation: "doubling",
		Cost:      "2M + 5S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/doubling/dbl-2009-alnr",
		Operation: "doubling",
		Cost:      "1M + 7S",
	},
	{
		ID:        "g1p/shortw/jacobian-0/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 8S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-0/doubling/dbl-1998-cmo-2",
		Operation: "doubling",
		Cost:      "3M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-0/doubling/dbl-1998-hnm",
		Operation: "doubling",
		Cost:      "3M + 6S + 1*half + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-0/doubling/dbl-1986-cc",
		Operation: "doubling",
		Cost:      "3M + 3S + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-0/doubling/dbl-1998-cmo",
		Operation: "doubling",
		Cost:      "3M + 3S + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-0/tripling/tpl-2007-bl",
		Operation: "tripling",
		Cost:      "5M + 10S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-0/tripling/tpl-2005-dim-2",
		Operation: "tripling",
		Cost:      "8M + 7S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-0/tripling/tpl-2005-dim",
		Operation: "tripling",
		Cost:      "9M + 5S + 1^3 + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-0/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 3M + 1S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/mmadd-2007-bl",
		Operation: "addition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/mmadd-2007-bl",
		Operation: "readdition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/zadd-2007-m",
		Operation: "addition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/zadd-2007-m",
		Operation: "readdition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/madd-2007-bl",
		Operation: "addition",
		Cost:      "7M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/madd-2007-bl",
		Operation: "readdition",
		Cost:      "7M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/madd-2004-hmv",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/madd-2004-hmv",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/madd-2008-g",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/madd-2008-g",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/madd",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/madd",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "11M + 5S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "10M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1998-cmo-2",
		Operation: "addition",
		Cost:      "12M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1998-cmo-2",
		Operation: "readdition",
		Cost:      "11M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-2001-b",
		Operation: "addition",
		Cost:      "12M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-2001-b",
		Operation: "readdition",
		Cost:      "11M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1998-hnm",
		Operation: "addition",
		Cost:      "12M + 4S + 1*half",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1998-hnm",
		Operation: "readdition",
		Cost:      "11M + 3S + 1*half",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1986-cc-2",
		Operation: "addition",
		Cost:      "8M + 6S + 2^3",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1986-cc-2",
		Operation: "readdition",
		Cost:      "8M + 5S + 1^3",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1986-cc",
		Operation: "addition",
		Cost:      "10M + 5S + 3^3",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1986-cc",
		Operation: "readdition",
		Cost:      "10M + 4S + 2^3",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1998-cmo",
		Operation: "addition",
		Cost:      "10M + 5S + 4^3",
	},
	{
		ID:        "g1p/shortw/jacobian-3/addition/add-1998-cmo",
		Operation: "readdition",
		Cost:      "10M + 4S + 3^3",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 5S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-2001-b",
		Operation: "doubling",
		Cost:      "3M + 5S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-2004-hmv",
		Operation: "doubling",
		Cost:      "4M + 4S + 1*half",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-1998-hnm-2",
		Operation: "doubling",
		Cost:      "4M + 4S + 1*half",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 8S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-1998-cmo-2",
		Operation: "doubling",
		Cost:      "3M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-1998-hnm",
		Operation: "doubling",
		Cost:      "3M + 6S + 1*half + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-1986-cc-2",
		Operation: "doubling",
		Cost:      "4M + 4S + 1^4",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-1986-cc",
		Operation: "doubling",
		Cost:      "3M + 3S + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-3/doubling/dbl-1998-cmo",
		Operation: "doubling",
		Cost:      "3M + 3S + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-3/tripling/tpl-2007-bl-2",
		Operation: "tripling",
		Cost:      "7M + 7S",
	},
	{
		ID:        "g1p/shortw/jacobian-3/tripling/tpl-2007-bl",
		Operation: "tripling",
		Cost:      "5M + 10S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-3/tripling/tpl-2005-dim-2",
		Operation: "tripling",
		Cost:      "8M + 7S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-3/tripling/tpl-2005-dim",
		Operation: "tripling",
		Cost:      "9M + 5S + 1^3 + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian-3/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 3M + 1S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/mmadd-2007-bl",
		Operation: "addition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/mmadd-2007-bl",
		Operation: "readdition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/zadd-2007-m",
		Operation: "addition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/zadd-2007-m",
		Operation: "readdition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/madd-2007-bl",
		Operation: "addition",
		Cost:      "7M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/madd-2007-bl",
		Operation: "readdition",
		Cost:      "7M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/madd-2004-hmv",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/madd-2004-hmv",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/madd-2008-g",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/madd-2008-g",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/madd",
		Operation: "addition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/madd",
		Operation: "readdition",
		Cost:      "8M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "11M + 5S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "10M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1998-cmo-2",
		Operation: "addition",
		Cost:      "12M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1998-cmo-2",
		Operation: "readdition",
		Cost:      "11M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-2001-b",
		Operation: "addition",
		Cost:      "12M + 4S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-2001-b",
		Operation: "readdition",
		Cost:      "11M + 3S",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1998-hnm",
		Operation: "addition",
		Cost:      "12M + 4S + 1*half",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1998-hnm",
		Operation: "readdition",
		Cost:      "11M + 3S + 1*half",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1986-cc-2",
		Operation: "addition",
		Cost:      "8M + 6S + 2^3",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1986-cc-2",
		Operation: "readdition",
		Cost:      "8M + 5S + 1^3",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1986-cc",
		Operation: "addition",
		Cost:      "10M + 5S + 3^3",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1986-cc",
		Operation: "readdition",
		Cost:      "10M + 4S + 2^3",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1998-cmo",
		Operation: "addition",
		Cost:      "10M + 5S + 4^3",
	},
	{
		ID:        "g1p/shortw/jacobian/addition/add-1998-cmo",
		Operation: "readdition",
		Cost:      "10M + 4S + 3^3",
	},
	{
		ID:        "g1p/shortw/jacobian/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 5S",
	},
	{
		ID:        "g1p/shortw/jacobian/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "1M + 8S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian/doubling/dbl-1998-cmo-2",
		Operation: "doubling",
		Cost:      "3M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian/doubling/dbl-1998-hnm",
		Operation: "doubling",
		Cost:      "3M + 6S + 1*half + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian/doubling/dbl-1986-cc",
		Operation: "doubling",
		Cost:      "3M + 3S + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian/doubling/dbl-1998-cmo",
		Operation: "doubling",
		Cost:      "3M + 3S + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian/tripling/tpl-2007-bl",
		Operation: "tripling",
		Cost:      "5M + 10S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian/tripling/tpl-2005-dim-2",
		Operation: "tripling",
		Cost:      "8M + 7S + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian/tripling/tpl-2005-dim",
		Operation: "tripling",
		Cost:      "9M + 5S + 1^3 + 2^4 + 1*a",
	},
	{
		ID:        "g1p/shortw/jacobian/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 3M + 1S",
	},
	{
		ID:        "g1p/shortw/modified/addition/mmadd-2009-bl",
		Operation: "addition",
		Cost:      "3M + 4S + 1*a",
	},
	{
		ID:        "g1p/shortw/modified/addition/mmadd-2009-bl",
		Operation: "readdition",
		Cost:      "3M + 4S + 1*a",
	},
	{
		ID:        "g1p/shortw/modified/addition/madd-2009-bl",
		Operation: "addition",
		Cost:      "7M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/modified/addition/madd-2009-bl",
		Operation: "readdition",
		Cost:      "7M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/modified/addition/add-2009-bl",
		Operation: "addition",
		Cost:      "11M + 7S + 1*a",
	},
	{
		ID:        "g1p/shortw/modified/addition/add-2009-bl",
		Operation: "readdition",
		Cost:      "10M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/modified/addition/add-1998-cmo-2",
		Operation: "addition",
		Cost:      "12M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/modified/addition/add-1998-cmo-2",
		Operation: "readdition",
		Cost:      "11M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/modified/doubling/mdbl-2009-bl",
		Operation: "doubling",
		Cost:      "2M + 5S",
	},
	{
		ID:        "g1p/shortw/modified/doubling/dbl-2009-bl",
		Operation: "doubling",
		Cost:      "3M + 5S",
	},
	{
		ID:        "g1p/shortw/modified/doubling/dbl-1998-cmo-2",
		Operation: "doubling",
		Cost:      "4M + 4S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/mmadd-1998-cmo",
		Operation: "addition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/mmadd-1998-cmo",
		Operation: "readdition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/madd-1998-cmo",
		Operation: "addition",
		Cost:      "9M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/madd-1998-cmo",
		Operation: "readdition",
		Cost:      "9M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-1998-cmo-2",
		Operation: "addition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-1998-cmo-2",
		Operation: "readdition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-2002-bj-2",
		Operation: "addition",
		Cost:      "13M + 3S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-2002-bj-2",
		Operation: "readdition",
		Cost:      "13M + 3S",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "11M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "11M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-2002-bj",
		Operation: "addition",
		Cost:      "12M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-2002-bj",
		Operation: "readdition",
		Cost:      "12M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-1986-cc",
		Operation: "addition",
		Cost:      "10M + 4S + 1^3",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-1986-cc",
		Operation: "readdition",
		Cost:      "10M + 4S + 1^3",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-1998-cmo",
		Operation: "addition",
		Cost:      "16M + 3S + 3^3",
	},
	{
		ID:        "g1p/shortw/projective-1/addition/add-1998-cmo",
		Operation: "readdition",
		Cost:      "16M + 3S + 3^3",
	},
	{
		ID:        "g1p/shortw/projective-1/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "3M + 5S",
	},
	{
		ID:        "g1p/shortw/projective-1/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "5M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-1/doubling/dbl-1998-cmo-2",
		Operation: "doubling",
		Cost:      "6M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-1/doubling/dbl-1998-cmo",
		Operation: "doubling",
		Cost:      "6M + 5S + 1^3 + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-1/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/mmadd-1998-cmo",
		Operation: "addition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/mmadd-1998-cmo",
		Operation: "readdition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/madd-1998-cmo",
		Operation: "addition",
		Cost:      "9M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/madd-1998-cmo",
		Operation: "readdition",
		Cost:      "9M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-1998-cmo-2",
		Operation: "addition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-1998-cmo-2",
		Operation: "readdition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "11M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "11M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-2002-bj",
		Operation: "addition",
		Cost:      "12M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-2002-bj",
		Operation: "readdition",
		Cost:      "12M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-1986-cc",
		Operation: "addition",
		Cost:      "10M + 4S + 1^3",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-1986-cc",
		Operation: "readdition",
		Cost:      "10M + 4S + 1^3",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-1998-cmo",
		Operation: "addition",
		Cost:      "16M + 3S + 3^3",
	},
	{
		ID:        "g1p/shortw/projective-3/addition/add-1998-cmo",
		Operation: "readdition",
		Cost:      "16M + 3S + 3^3",
	},
	{
		ID:        "g1p/shortw/projective-3/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "3M + 5S",
	},
	{
		ID:        "g1p/shortw/projective-3/doubling/dbl-2007-bl-2",
		Operation: "doubling",
		Cost:      "7M + 3S",
	},
	{
		ID:        "g1p/shortw/projective-3/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "5M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-3/doubling/dbl-1998-cmo-2",
		Operation: "doubling",
		Cost:      "6M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-3/doubling/dbl-1998-cmo",
		Operation: "doubling",
		Cost:      "6M + 5S + 1^3 + 1*a",
	},
	{
		ID:        "g1p/shortw/projective-3/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g1p/shortw/projective/addition/mmadd-1998-cmo",
		Operation: "addition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/projective/addition/mmadd-1998-cmo",
		Operation: "readdition",
		Cost:      "5M + 2S",
	},
	{
		ID:        "g1p/shortw/projective/addition/madd-1998-cmo",
		Operation: "addition",
		Cost:      "9M + 2S",
	},
	{
		ID:        "g1p/shortw/projective/addition/madd-1998-cmo",
		Operation: "readdition",
		Cost:      "9M + 2S",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-1998-cmo-2",
		Operation: "addition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-1998-cmo-2",
		Operation: "readdition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-2007-bl",
		Operation: "addition",
		Cost:      "11M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-2007-bl",
		Operation: "readdition",
		Cost:      "11M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-2002-bj",
		Operation: "addition",
		Cost:      "12M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-2002-bj",
		Operation: "readdition",
		Cost:      "12M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-1986-cc",
		Operation: "addition",
		Cost:      "10M + 4S + 1^3",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-1986-cc",
		Operation: "readdition",
		Cost:      "10M + 4S + 1^3",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-1998-cmo",
		Operation: "addition",
		Cost:      "16M + 3S + 3^3",
	},
	{
		ID:        "g1p/shortw/projective/addition/add-1998-cmo",
		Operation: "readdition",
		Cost:      "16M + 3S + 3^3",
	},
	{
		ID:        "g1p/shortw/projective/doubling/mdbl-2007-bl",
		Operation: "doubling",
		Cost:      "3M + 5S",
	},
	{
		ID:        "g1p/shortw/projective/doubling/dbl-2007-bl",
		Operation: "doubling",
		Cost:      "5M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective/doubling/dbl-1998-cmo-2",
		Operation: "doubling",
		Cost:      "6M + 5S + 1*a",
	},
	{
		ID:        "g1p/shortw/projective/doubling/dbl-1998-cmo",
		Operation: "doubling",
		Cost:      "6M + 5S + 1^3 + 1*a",
	},
	{
		ID:        "g1p/shortw/projective/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 2M",
	},
	{
		ID:        "g1p/shortw/w12-0/addition/madd-2010-cln",
		Operation: "addition",
		Cost:      "8M + 5S",
	},
	{
		ID:        "g1p/shortw/w12-0/addition/madd-2010-cln",
		Operation: "readdition",
		Cost:      "8M + 5S",
	},
	{
		ID:        "g1p/shortw/w12-0/addition/add-2010-cln",
		Operation: "addition",
		Cost:      "10M + 7S",
	},
	{
		ID:        "g1p/shortw/w12-0/addition/add-2010-cln",
		Operation: "readdition",
		Cost:      "10M + 6S",
	},
	{
		ID:        "g1p/shortw/w12-0/doubling/dbl-2010-cln",
		Operation: "doubling",
		Cost:      "1M + 6S + 1*a",
	},
	{
		ID:        "g1p/shortw/xyzz-3/addition/mmadd-2008-s",
		Operation: "addition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/addition/mmadd-2008-s",
		Operation: "readdition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/addition/madd-2008-s",
		Operation: "addition",
		Cost:      "8M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/addition/madd-2008-s",
		Operation: "readdition",
		Cost:      "8M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/addition/add-2008-s",
		Operation: "addition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/addition/add-2008-s",
		Operation: "readdition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/doubling/mdbl-2008-s-1",
		Operation: "doubling",
		Cost:      "4M + 3S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/doubling/mdbl-2008-s-2",
		Operation: "doubling",
		Cost:      "4M + 3S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/doubling/dbl-2008-s-2",
		Operation: "doubling",
		Cost:      "7M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz-3/doubling/dbl-2008-s-1",
		Operation: "doubling",
		Cost:      "6M + 4S + 1*a",
	},
	{
		ID:        "g1p/shortw/xyzz-3/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 3M + 1S",
	},
	{
		ID:        "g1p/shortw/xyzz/addition/mmadd-2008-s",
		Operation: "addition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz/addition/mmadd-2008-s",
		Operation: "readdition",
		Cost:      "4M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz/addition/madd-2008-s",
		Operation: "addition",
		Cost:      "8M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz/addition/madd-2008-s",
		Operation: "readdition",
		Cost:      "8M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz/addition/add-2008-s",
		Operation: "addition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz/addition/add-2008-s",
		Operation: "readdition",
		Cost:      "12M + 2S",
	},
	{
		ID:        "g1p/shortw/xyzz/doubling/mdbl-2008-s-1",
		Operation: "doubling",
		Cost:      "4M + 3S",
	},
	{
		ID:        "g1p/shortw/xyzz/doubling/dbl-2008-s-1",
		Operation: "doubling",
		Cost:      "6M + 4S + 1*a",
	},
	{
		ID:        "g1p/shortw/xyzz/scaling/z",
		Operation: "scaling",
		Cost:      "1I + 3M + 1S",
	},
	{
		ID:        "g1p/shortw/xz/doubling/dbl-2002-bj-3",
		Operation: "doubling",
		Cost:      "2M + 5S + 1*b2 + 1*a + 1*b4",
	},
	{
		ID:        "g1p/shortw/xz/doubling/dbl-2002-bj-2",
		Operation: "doubling",
		Cost:      "3M + 4S + 1*b2 + 1*a",
	},
	{
		ID:        "g1p/shortw/xz/doubling/dbl-2002-it-2",
		Operation: "doubling",
		Cost:      "4M + 3S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/doubling/dbl-2002-it",
		Operation: "doubling",
		Cost:      "3M + 5S + 1^3 + 1^4 + 2*a + 2*b",
	},
	{
		ID:        "g1p/shortw/xz/doubling/dbl-2002-bj",
		Operation: "doubling",
		Cost:      "3M + 4S + 3^3 + 2*a + 2*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/mdadd-2002-bj-2",
		Operation: "diffadd",
		Cost:      "6M + 2S + 1*a + 1*b4",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/mdadd-2002-it-3",
		Operation: "diffadd",
		Cost:      "6M + 2S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/mdadd-2002-it-4",
		Operation: "diffadd",
		Cost:      "6M + 2S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/dadd-2002-it-3",
		Operation: "diffadd",
		Cost:      "7M + 2S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/dadd-2002-it-4",
		Operation: "diffadd",
		Cost:      "8M + 2S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/mdadd-2002-bj",
		Operation: "diffadd",
		Cost:      "9M + 2S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/mdadd-2002-it",
		Operation: "diffadd",
		Cost:      "9M + 2S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/mdadd-2002-it-2",
		Operation: "diffadd",
		Cost:      "9M + 3S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/dadd-2002-it",
		Operation: "diffadd",
		Cost:      "10M + 2S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/diffadd/dadd-2002-it-2",
		Operation: "diffadd",
		Cost:      "11M + 3S + 1*a + 1*b",
	},
	{
		ID:        "g1p/shortw/xz/ladder/mladd-2002-bj-3",
		Operation: "ladder",
		Cost:      "8M + 7S + 1*b2 + 2*a + 2*b4",
	},
	{
		ID:        "g1p/shortw/xz/ladder/mladd-2002-it-3",
		Operation: "ladder",
		Cost:      "8M + 7S + 2*a + 3*b4",
	},
	{
		ID:        "g1p/shortw/xz/ladder/mladd-2002-it-4",
		Operation: "ladder",
		Cost:      "8M + 7S + 2*a + 3*b4",
	},
	{
		ID:        "g1p/shortw/xz/ladder/mladd-2002-it-5",
		Operation: "ladder",
		Cost:      "8M + 7S + 2*a + 3*b4",
	},
	{
		ID:        "g1p/shortw/xz/ladder/mladd-2002-bj-2",
		Operation: "ladder",
		Cost:      "9M + 6S + 1*b2 + 2*a + 1*b4",
	},
	{
		ID:        "g1p/shortw/xz/ladder/ladd-2002-it-3",
		Operation: "ladder",
		Cost:      "9M + 7S + 2*a + 3*b4",
	},
	{
		ID:        "g1p/shortw/xz/ladder/ladd-2002-it-4",
		Operation: "ladder",
		Cost:      "10M + 7S + 2*a + 3*b4",
	},
	{
		ID:        "g1p/shortw/xz/ladder/mladd-2002-it",
		Operation: "ladder",
		Cost:      "12M + 7S + 1^3 + 1^4 + 3*a + 3*b",
	},
	{
		ID:        "g1p/shortw/xz/ladder/mladd-2002-it-2",
		Operation: "ladder",
		Cost:      "12M + 8S + 1^3 + 1^4 + 3*a + 3*b",
	},
	{
		ID:        "g1p/shortw/xz/ladder/ladd-2002-it",
		Operation: "ladder",
		Cost:      "13M + 7S + 1^3 + 1^4 + 3*a + 3*b",
	},
	{
		ID:        "g1p/shortw/xz/ladder/ladd-2002-it-2",
		Operation: "ladder",
		Cost:      "14M + 8S + 1^3 + 1^4 + 3*a + 3*b",
	},
	{
		ID:        "g1p/shortw/xz/ladder/mladd-2002-bj",
		Operation: "ladder",
		Cost:      "12M + 6S + 3^3 + 3*a + 3*b",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/mmadd-2008-hwcd-4",
		Operation: "addition",
		Cost:      "6M",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/mmadd-2008-hwcd-4",
		Operation: "readdition",
		Cost:      "6M",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/mmadd-2008-hwcd-3",
		Operation: "addition",
		Cost:      "6M + 1S + 1*k",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/mmadd-2008-hwcd-3",
		Operation: "readdition",
		Cost:      "6M + 1S",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/mmadd-2008-hwcd-2",
		Operation: "addition",
		Cost:      "7M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/mmadd-2008-hwcd-2",
		Operation: "readdition",
		Cost:      "7M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/madd-2008-hwcd-3",
		Operation: "addition",
		Cost:      "7M + 1*k",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/madd-2008-hwcd-3",
		Operation: "readdition",
		Cost:      "7M",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/madd-2008-hwcd-4",
		Operation: "addition",
		Cost:      "7M",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/madd-2008-hwcd-4",
		Operation: "readdition",
		Cost:      "7M",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/mmadd-2008-hwcd",
		Operation: "addition",
		Cost:      "7M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/mmadd-2008-hwcd",
		Operation: "readdition",
		Cost:      "7M + 1S + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/madd-2008-hwcd-2",
		Operation: "addition",
		Cost:      "8M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/madd-2008-hwcd-2",
		Operation: "readdition",
		Cost:      "8M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/madd-2008-hwcd",
		Operation: "addition",
		Cost:      "8M + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/madd-2008-hwcd",
		Operation: "readdition",
		Cost:      "8M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/add-2008-hwcd-3",
		Operation: "addition",
		Cost:      "8M + 1*k",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/add-2008-hwcd-3",
		Operation: "readdition",
		Cost:      "8M",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/add-2008-hwcd-4",
		Operation: "addition",
		Cost:      "8M",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/add-2008-hwcd-4",
		Operation: "readdition",
		Cost:      "8M",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/add-2008-hwcd-2",
		Operation: "addition",
		Cost:      "9M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/add-2008-hwcd-2",
		Operation: "readdition",
		Cost:      "9M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/add-2008-hwcd",
		Operation: "addition",
		Cost:      "9M + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/extended-1/addition/add-2008-hwcd",
		Operation: "readdition",
		Cost:      "9M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/doubling/mdbl-2008-hwcd",
		Operation: "doubling",
		Cost:      "3M + 4S + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/doubling/dbl-2008-hwcd",
		Operation: "doubling",
		Cost:      "4M + 4S + 1*a",
	},
	{
		ID:        "g1p/twisted/extended-1/tripling/tpl-2015-c",
		Operation: "tripling",
		Cost:      "11M + 3S + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/mmadd-2008-hwcd-2",
		Operation: "addition",
		Cost:      "7M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/mmadd-2008-hwcd-2",
		Operation: "readdition",
		Cost:      "7M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/mmadd-2008-hwcd",
		Operation: "addition",
		Cost:      "7M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/extended/addition/mmadd-2008-hwcd",
		Operation: "readdition",
		Cost:      "7M + 1S + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/madd-2008-hwcd-2",
		Operation: "addition",
		Cost:      "8M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/madd-2008-hwcd-2",
		Operation: "readdition",
		Cost:      "8M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/madd-2008-hwcd",
		Operation: "addition",
		Cost:      "8M + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/extended/addition/madd-2008-hwcd",
		Operation: "readdition",
		Cost:      "8M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/add-2008-hwcd-2",
		Operation: "addition",
		Cost:      "9M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/add-2008-hwcd-2",
		Operation: "readdition",
		Cost:      "9M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/addition/add-2008-hwcd",
		Operation: "addition",
		Cost:      "9M + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/extended/addition/add-2008-hwcd",
		Operation: "readdition",
		Cost:      "9M + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/doubling/mdbl-2008-hwcd",
		Operation: "doubling",
		Cost:      "3M + 4S + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/doubling/dbl-2008-hwcd",
		Operation: "doubling",
		Cost:      "4M + 4S + 1*a",
	},
	{
		ID:        "g1p/twisted/extended/tripling/tpl-2015-c",
		Operation: "tripling",
		Cost:      "11M + 3S + 1*a",
	},
	{
		ID:        "g1p/twisted/inverted/addition/mmadd-2008-bbjlp",
		Operation: "addition",
		Cost:      "7M + 1*a",
	},
	{
		ID:        "g1p/twisted/inverted/addition/mmadd-2008-bbjlp",
		Operation: "readdition",
		Cost:      "7M + 1*a",
	},
	{
		ID:        "g1p/twisted/inverted/addition/madd-2008-bbjlp",
		Operation: "addition",
		Cost:      "8M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/inverted/addition/madd-2008-bbjlp",
		Operation: "readdition",
		Cost:      "8M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/inverted/addition/add-2008-bbjlp",
		Operation: "addition",
		Cost:      "9M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/inverted/addition/add-2008-bbjlp",
		Operation: "readdition",
		Cost:      "9M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/inverted/doubling/mdbl-2008-bbjlp",
		Operation: "doubling",
		Cost:      "3M + 3S + 1*a",
	},
	{
		ID:        "g1p/twisted/inverted/doubling/dbl-2008-bbjlp",
		Operation: "doubling",
		Cost:      "3M + 4S + 1*a + 1*d2",
	},
	{
		ID:        "g1p/twisted/projective/addition/mmadd-2008-bbjlp",
		Operation: "addition",
		Cost:      "6M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/projective/addition/mmadd-2008-bbjlp",
		Operation: "readdition",
		Cost:      "6M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/projective/addition/madd-2008-bbjlp",
		Operation: "addition",
		Cost:      "9M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/projective/addition/madd-2008-bbjlp",
		Operation: "readdition",
		Cost:      "9M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/projective/addition/add-2008-bbjlp",
		Operation: "addition",
		Cost:      "10M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/projective/addition/add-2008-bbjlp",
		Operation: "readdition",
		Cost:      "10M + 1S + 1*a + 1*d",
	},
	{
		ID:        "g1p/twisted/projective/doubling/mdbl-2008-bbjlp",
		Operation: "doubling",
		Cost:      "2M + 4S + 1*a",
	},
	{
		ID:        "g1p/twisted/projective/doubling/dbl-2008-bbjlp",
		Operation: "doubling",
		Cost:      "3M + 4S + 1*a",
	},
	{
		ID:        "g1p/twisted/projective/tripling/tpl-2015-c",
		Operation: "tripling",
		Cost:      "9M + 3S + 1*a",
	},
	{
		ID:        "g1p/twistedhessian/projective/addition/add-2009-bkl",
		Operation: "addition",
		Cost:      "12M + 1*a",
	},
	{
		ID:        "g1p/twistedhessian/projective/addition/add-2009-bkl",
		Operation: "readdition",
		Cost:      "12M + 1*a",
	},
	{
		ID:        "g1p/twistedhessian/projective/doubling/dbl-2012-c",
		Operation: "doubling",
		Cost:      "7M + 1S + 1*minustwo + 1*d",
	},
	{
		ID:        "g1p/twistedhessian/projective/doubling/dbl-2009-bkl-3",
		Operation: "doubling",
		Cost:      "8M + 1*i + 1*minustwo + 1*2d",
	},
	{
		ID:        "g1p/twistedhessian/projective/doubling/dbl-2009-bkl-2",
		Operation: "doubling",
		Cost:      "6M + 3S + 1*a",
	},
	{
		ID:        "g1p/twistedhessian/projective/doubling/dbl-2009-bkl",
		Operation: "doubling",
		Cost:      "3M + 3^3 + 1*a",
	},
	{
		ID:        "g1p/twistedhessian/projective/tripling/tpl-2009-bkl",
		Operation: "tripling",
		Cost:      "8M + 6S + 1*a + 1*recipd",
	},
}
